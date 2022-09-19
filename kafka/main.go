package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
)

var wg sync.WaitGroup

var (
	config  *sarama.Config
	host    string
	port    int
	topic   string
	produce bool
	help    bool
)

func init() {
	flag.StringVar(&host, "host", "localhost", "kafka server host")
	flag.IntVar(&port, "port", 9092, "kafka server port")
	flag.StringVar(&topic, "topic", "test.topic", "topic to produce and consume")
	flag.BoolVar(&help, "h", false, "show help info")
	flag.Parse()
}

func main() {
	if help {
		flag.PrintDefaults()
		return
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config = sarama.NewConfig()
	//ack应答机制
	config.Producer.RequiredAcks = sarama.WaitForAll
	//发送分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//回复确认
	config.Producer.Return.Successes = true
	addr := fmt.Sprintf("%s:%d", host, port)

	if produce {
		push(addr, topic)
	}
	consume(addr, topic)
}

func push(addr, topic string) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder("a test message [hello world]")
	//连接kafka
	client, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		log.Fatalf("new producer failed: %v", err)
	}
	defer client.Close()
	//发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		log.Printf("send msg failed: %v", err)
		return
	}
	log.Printf("pid:%v offset:%v", pid, offset)
}

func consume(addr, topic string) {
	consumer, err := sarama.NewConsumer([]string{addr}, nil)
	if err != nil {
		log.Println("fail to start consumer", err)
	}
	// get partitions
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Println("fail to get list of partition, err:", err)
	}
	log.Println(partitionList)
	for p := range partitionList {
		// create consumer for each partiton
		pc, err := consumer.ConsumePartition(topic, int32(p), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d, err:%v\n", p, err)
		}
		defer pc.AsyncClose()
		wg.Add(1)
		// consume
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				log.Printf("partition:%d Offset:%d Key:%v Value:%s\n",
					msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
			wg.Done()
		}(pc)
	}
	wg.Wait()
}
