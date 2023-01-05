package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/brinkpku/training_center/go-redis/queue"
	"github.com/go-redis/redis/v8"
)

func main() {
	mainTest()
	return
	ctx := context.Background()
	q, err := queue.NewCircleQueue(&queue.Config{
		DSN: "redis://:Hp2HZuMj@10.10.18.109:61234/3",
	})
	if err != nil {
		panic(err)
	}
	q1, err := queue.NewCircleQueue(&queue.Config{
		DSN: "redis://:Hp2HZuMj@10.10.18.109:61234/3",
	})
	if err != nil {
		panic(err)
	}
	q.Init(ctx, "test:q")
	q1.Init(ctx, "test:q")
	for i := 1; i < 10; i++ {
		if err := q.Push(ctx, i); err != nil {
			panic(err)
		}
	}
	l, err := q.Length(ctx)
	if err != nil {
		panic(err)
	}
	println("queue length:", l)
	wg := &sync.WaitGroup{}
	for i := 0; i < 2; i++ { // mock multi-client
		wg.Add(1)
		go func(id int) {
			ch := make(chan string, 1)
			go func() {
				err := q.Traverse(ctx, ch)
				if err != nil {
					panic(err)
				}
			}()
			for ele := range ch {
				println("consumer", id, "consumed", ele)
			}
			wg.Done()
		}(i)
	}

	println("wait client running..")
	wg.Wait()
	q.Flush(ctx)
}

type multiPartUploadDesc struct {
	sync.Mutex

	JobID   string
	MaxSize int
	MaxNum  int
	Packs   []*pack
}

type pack struct {
	test string
	num  int
}

// redis marshaller
func (mpd *multiPartUploadDesc) MarshalBinary() ([]byte, error) {
	if mpd == nil {
		mpd = new(multiPartUploadDesc)
	}
	return json.Marshal(mpd)
}

// redis scan unmarshaller
func (mpd *multiPartUploadDesc) UnmarshalBinary(data []byte) (err error) {
	if mpd == nil {
		mpd = new(multiPartUploadDesc)
	}
	return json.Unmarshal(data, mpd)
}

func mainTest() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.10.18.109:61234",
		Password: "Hp2HZuMj", // no password set
		DB:       6,          // use default DB
	})
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	nt := &multiPartUploadDesc{}
	if err = rdb.Get(ctx, "multipartUploadDesc:81").Scan(nt); err != nil {
		panic(err)
	}
	fmt.Println(nt)
	fmt.Println(nt.Packs)
	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
