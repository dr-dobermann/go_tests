package main

import (
	"context"
	"fmt"

	kafka "github.com/segmentio/kafka-go"
)

func main() {

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     "invoices",
		Partition: 0,
		MinBytes:  10e3, // 10Kb
		MaxBytes:  10e6, // 10Mb
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())

		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

}
