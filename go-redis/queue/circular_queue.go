package queue

/*
thread-safe circle queue based on redis list.
*/

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// FRONT mark the start of circle queue
	FRONT = "FRONT"
)

type CircularQueue interface {
	Flush(context.Context) error
	Init(ctx context.Context, queueName string) error
	IsFront(context.Context) (bool, error)
	Length(context.Context) (int64, error)
	Lock(context.Context) (bool, error)
	Push(context.Context, ...interface{}) error
	Rotate(context.Context) (string, error)
	Traverse(context.Context, chan<- string) error
	Unlock(context.Context) error
}

type Config struct {
	QueueName    string
	DSN          string        `toml:"dsn" json:"dsn"`
	MaxIdle      int           `toml:"db_conn_pool_max_idle" json:"db_conn_pool_max_idle"`         // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int           `toml:"db_conn_pool_max_open" json:"db_conn_pool_max_open"`         // <= 0 means unlimited
	MaxLifetime  time.Duration `toml:"db_conn_pool_max_lifetime" json:"db_conn_pool_max_lifetime"` // maximum amount of time a connection may be reused
	PrintDBStats bool          `toml:"print_db_stats" json:"print_db_stats"`
}

func NewCircleQueue(config *Config) (CircularQueue, error) {
	opt, err := redis.ParseURL(config.DSN)
	if err != nil {
		return nil, err
	}
	cli := redis.NewClient(opt)
	return &circleQueueImpl{
		cli: cli,
	}, nil
}

type lock struct {
	key   string
	value string
}

type circleQueueImpl struct {
	cli  *redis.Client
	lock *lock

	cfg          *Config
	queueName    string
	isIterateKey string
}

// Init
func (cq *circleQueueImpl) Init(ctx context.Context, name string) (err error) {
	cq.queueName = name
	cq.isIterateKey = fmt.Sprintf("%s:is_iterate", name)
	cq.lock = &lock{
		key: fmt.Sprintf("%s:lock", name),
	}
	locked, err := cq.Lock(ctx)
	if err != nil {
		return
	}
	if locked {
		defer cq.Unlock(ctx)
	}
	exist, err := cq.cli.Exists(ctx, name).Result()
	if err != nil {
		return
	}
	if exist == 0 { // TODO use tx
		err = cq.Push(ctx, FRONT)
		_, err = cq.cli.Set(ctx, cq.isIterateKey, "false", 0).Result()
	}
	return
}

// Lock lock circle queue
func (cq *circleQueueImpl) Lock(ctx context.Context) (bool, error) {
	return cq.lockKeyWithTTL(ctx, cq.lock.key, time.Second*2)
}

// lock circle queue with ttl
func (cq *circleQueueImpl) lockKeyWithTTL(ctx context.Context, key string, dur time.Duration) (bool, error) {
	cq.lock.value = fmt.Sprint(time.Now().Add(dur).Unix())
	return cq.cli.SetNX(ctx, cq.lock.key, cq.lock.value, dur).Result()
}

// Unlock
func (cq *circleQueueImpl) Unlock(ctx context.Context) (err error) {
	_, err = cq.cli.Del(ctx, cq.lock.key).Result()
	return
}

func (cq *circleQueueImpl) Length(ctx context.Context) (int64, error) {
	return cq.cli.LLen(ctx, cq.queueName).Result()
}

func (cq *circleQueueImpl) front(ctx context.Context) (string, error) {
	return cq.cli.LIndex(ctx, cq.queueName, 0).Result()
}

func (cq *circleQueueImpl) IsFront(ctx context.Context) (isFront bool, err error) {
	ele, err := cq.front(ctx)
	if err != nil {
		return
	}
	isFront = (ele == FRONT)
	return
}

func (cq *circleQueueImpl) Push(ctx context.Context, elements ...interface{}) (err error) {
	_, err = cq.cli.RPush(ctx, cq.queueName, elements...).Result()
	return
}

func (cq *circleQueueImpl) Rotate(ctx context.Context) (ele string, err error) {
	ele, err = cq.cli.LMove(ctx, cq.queueName, cq.queueName, "LEFT", "RIGHT").Result()
	return
}

// Flush clean all elements of cq.
func (cq *circleQueueImpl) Flush(ctx context.Context) (err error) {
	if _, err = cq.cli.Del(ctx, cq.queueName).Result(); err != nil { // TODO use transaction
		return
	}
	_, err = cq.cli.Del(ctx, cq.lock.key).Result()
	return
}

// Traverse traverse queue one round, thread-safe.
func (cq *circleQueueImpl) Traverse(ctx context.Context, ch chan<- string) error {
	for {
		locked, err := cq.lockKeyWithTTL(ctx, cq.isIterateKey, time.Minute)
		if err != nil {
			return err
		}
		if locked {
			isFront, err := cq.IsFront(ctx)
			if err != nil {
				return err
			}
			if isFront {
				_, err := cq.Rotate(ctx) // rotate FRONT
				if err != nil {
					return err
				}
			}
		}
		iterate, err := cq.cli.Get(ctx, cq.isIterateKey).Result()
		if err != nil {
			return err
		}
		isIterate, err := strconv.ParseBool(iterate)
		if err != nil {
			return err
		}
		if isIterate {
			isFront, err := cq.IsFront(ctx)
			if err != nil {
				return err
			}
			if isFront {
				cq.cli.Set(ctx, cq.isIterateKey, "false", 0)
				close(ch)
				return nil
			}
			ele, err := cq.Rotate(ctx)
			if err != nil {
				return err
			}
			ch <- ele
		} else {
			locked, err := cq.Lock(ctx)
			if err != nil {
				return err
			}
			if locked {
				println("locked")
				// not release lock, wait one job time
				// defer cq.Unlock(ctx)
			} else { // cant get lock, retry
				time.Sleep(time.Second)
				println("blocked")
				continue
			}
			isFront, err := cq.IsFront(ctx)
			if err != nil {
				return err
			}
			if isFront {
				_, err := cq.Rotate(ctx) // rotate FRONT
				if err != nil {
					return err
				}
			}
			_, err = cq.cli.Set(ctx, cq.isIterateKey, "true", 0).Result()
			if err != nil {
				return err
			}
		}
	}
}
