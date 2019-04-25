package main

import (
	"context"
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"
	uuid "github.com/satori/go.uuid"
)

func main() {

	start := time.Now()
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "invoices",
		Balancer: &kafka.LeastBytes{},
	})

	ch := make(chan int, 3)
	quit := make(chan int)

	var i int

	for i = 0; i < 3; i++ {
		go GenerateMessages(w, 100, ch, quit)
	}

	i = 3
	for i > 0 {
		select {
			case n := <- ch:
				fmt.Println(n)
			case <-quit:
				i-- 
		}
	}

	w.Close()
	t := time.Now()
	fmt.Println(t.Sub(start).Seconds())
}

// GenerateMessages sends num generated messages into invoice topic.
// After every message it sends its number into ch channel and sens 0 into quit channel on finish
func GenerateMessages(w *kafka.Writer, num int, ch chan int, quit chan int) {
	for i := 0; i < num; i++ {
		uid := uuid.NewV4()
		w.WriteMessages(
			context.Background(), 
			kafka.Message{
				Key:   []byte("invoice"),
				Value: []byte(uid.String()),
			},)
		ch <- i
	}
	quit <- 0
}