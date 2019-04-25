package main

import (
	"context"
	"fmt"
	"flag"
	"strings"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func main() {

	var (
		bs = flag.String("brokers", "localhost:9092", "Enter broker lists, separated by comma (Example: \"kafka1:9092,kafka2:9092\")")
		t = flag.String("topic", "echo", "Topic name to subscribe")
		p = flag.Int("partition", 0, "Kafka partition number")
		o = flag.Int("offset", 0, "Topic initial offset for reading")
		short = flag.Bool("short", true, "Show short message value")
		bb []string
	)

	flag.Parse()
	if *p < 0 {
		log.Fatalf("Invalid partiton number %v. Should be equal or highter than zero!", *p)
	}
	if *o < 0 {
		log.Fatalf("Invalid offset number %v. Should be greater than zero!", *o)
	}
	bb = strings.Split(*bs, ",")
	if len(bb) == 0 {
		log.Fatal("Should be set at least one broker!")
	}
	for i, b := range bb {
		bb[i] = strings.TrimSpace(b)
	}

	fmt.Printf("Trying to subscribe on topic: [%s]\nhosted on kafka broker(s): %v\nFrom offset: %d\n", *t, bb, *o)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   bb,
		Topic:     *t,
		Partition: *p,
		MinBytes:  10e3, // 10Kb
		MaxBytes:  10e6, // 10Mb
	})
	defer r.Close()

	r.SetOffset(int64(*o))
	for {
		m, err := r.ReadMessage(context.Background())

		if err != nil {
			break
		}
		msg := string(m.Value)
		if *short && len(msg) > 50 {
			msg = msg[:50]
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), msg)
	}

}
