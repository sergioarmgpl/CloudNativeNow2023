package main

import (
    "fmt"
    "os"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
    topic := "matches"
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "my-cluster-kafka-0.my-cluster-kafka-brokers.kafka.svc:9092"})
    if err != nil {
        fmt.Printf("Failed to create producer: %s", err)
        os.Exit(1)
    }

    p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Key: []byte("data"),
        Value: []byte("{\"data\":\"mydata\"}"),
    }, nil)

    // Wait for all messages to be delivered
    p.Flush(15 * 1000)
    p.Close()
}