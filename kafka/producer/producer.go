package main

/*
Docker containers from confluent Inc are using in this project as followed

docker network create kafka

docker run --network kafka -d --name=zookeeper -e ZOOKEEPER_CLIENT_PORT=2181 confluentinc/cp-zookeeper:latest

docker run --network kafka -d -p 9092:9092 --name=kafka -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092 -e KAFKA_OFFSET_TOPIC_REPLICATION_TOPIC=1 confluentinc/cp-kafka:latest

*/

import (
	"context"
	"fmt"
	"time"
	"flag"
	"log"
	"strings"
	"io/ioutil"

	kafka "github.com/segmentio/kafka-go"
	uuid "github.com/satori/go.uuid"
)

func main() {

	var (
		bs = flag.String("brokers", "localhost:9092", "Enter broker lists, separated by comma (Example: \"kafka1:9092,kafka2:9092\")")
		t = flag.String("topic", "echo", "Topic name to subscribe")
		s = flag.Int("streams", 3, "Number of parallel streams to generate messages")
		n = flag.Int("num", 10, "Number of generated messages in every stream")
		fn = flag.String("fname", "", "File name with a single data package")
		bb []string
	)

	flag.Parse()
	if *s < 1 {
		log.Fatalf("Invalid partiton number %v. Should be highter than zero", *s)
	}
	if *n < 1 {
		log.Fatal("Number of generated messages per stream should be higher than zero!", *n)
	}
	bb = strings.Split(*bs, ",")
	if len(bb) == 0 {
		log.Fatal("Should be set at least one broker!")
	}
	for i, b := range bb {
		bb[i] = strings.TrimSpace(b)
	}
	
	var (
		dat []byte
		err error
	)

	if len(*fn) != 0 {
		fmt.Printf("Trying to read data packet from file %s...\n", *fn)
		dat, err = ioutil.ReadFile(*fn)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("Data package readed:\nBEGIN\n----------\n%v\n----------\nEND\n", string(dat))
	} else {
		fmt.Printf("No filename for data packet was given. UUIDs will be used instead...\n")
	}

	fmt.Printf("Generating %d messages in %d stream(s) into topic: [%s]\nhosted on kafka broker(s): %v\n", *n, *s, *t, bb)

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  bb,
		Topic:    *t,
		Balancer: &kafka.LeastBytes{},
	})
	
	defer w.Close()

	ch := make(chan int, 3)
	quit := make(chan int)

	var i int
	start := time.Now()

	for i = 0; i < *s; i++ {
		go GenerateMessages(w, *n, ch, quit, dat)
	}

	tnum := *s * *n
	for *s > 0 {
		select {
			case i := <- ch:
				fmt.Println(i)
			case <-quit:
				*s-- 
		}
	}

	et := time.Now().Sub(start).Seconds()
	fmt.Printf("Elapsed time: %f seconds.\nMessages per seconds: %f\n", et, float64(tnum)/et)
}

// GenerateMessages sends the num of generated messages into invoice topic.
// After posting every message into kafka it sends its number into the ch channel.
// At the end it sends 0 into the quit channel.
func GenerateMessages(w *kafka.Writer, num int, ch chan int, quit chan int, data []byte) {
	var (
		key []byte
		value []byte
	)

	if len(data) == 0 {
		key = []byte("id")
	} else {
		key = []byte("data")
		value = data
	}
	for i := 0; i < num; i++ {
		if len(data) == 0 {
			value = []byte(uuid.NewV4().String())
			log.Println("UUID used...")
		}
		w.WriteMessages(
			context.Background(), 
			kafka.Message{
				Key:   key,
				Value: value,
			},)
		ch <- i
	}
	quit <- 0
}
