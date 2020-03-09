package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/adjust/rmq"
	"github.com/ibamaulana/golang-kumparan/request/news"
)

const (
	unackedLimit = 1000
	numConsumers = 10
	batchSize    = 1000
)

func main() {
	connection := rmq.OpenConnection("consumer", "tcp", "localhost:6379", 2)
	// queue := connection.OpenQueue("things")
	// queue.StartConsuming(unackedLimit, 500*time.Millisecond)

	queuetest := connection.OpenQueue("test")
	queuetest.StartConsuming(10, time.Second)
	queuetest.AddConsumer("consumer", NewConsumer(1))
	// for i := 0; i < numConsumers; i++ {
	// 	name := fmt.Sprintf("consumer %d", i)
	// 	queue.AddConsumer(name, NewConsumer(i))
	// }
	select {}
}

type Consumer struct {
	name   string
	count  int
	before time.Time
}

func NewConsumer(tag int) *Consumer {
	return &Consumer{
		name:   fmt.Sprintf("consumer%d", tag),
		count:  0,
		before: time.Now(),
	}
}

func (consumer *Consumer) Consume(delivery rmq.Delivery) {
	// consumer.count++
	// if consumer.count%batchSize == 0 {
	// 	duration := time.Now().Sub(consumer.before)
	// 	consumer.before = time.Now()
	// 	perSecond := time.Second / (duration / batchSize)
	// 	log.Printf("%s consumed %d %s %d", consumer.name, consumer.count, delivery.Payload(), perSecond)
	// }
	// time.Sleep(time.Millisecond)
	// if consumer.count%batchSize == 0 {
	// 	delivery.Reject()
	// } else {
	// 	delivery.Ack()
	// }
	var err error
	var task news.CreateRequest
	if err = json.Unmarshal([]byte(delivery.Payload()), &task); err != nil {
		// handle error
		log.Printf(delivery.Payload())
		delivery.Reject()
		return
	}

	// perform task
	log.Printf("performing task %s", task)
	delivery.Ack()
}
