package main

import (
	"github.com/segmentio/kafka-go"
	"os"
)

const (
	topicName = "test-topic"
	groupName = "test-group"
)

// getKafkaAddress returns Kafka broker address from environment or default
func getKafkaAddress() string {
	if addr := os.Getenv("KAFKA_ADDR"); addr != "" {
		return addr
	}
	return "127.0.0.1:9092" // Default Kafka address
}

// initProducer creates a new Kafka producer configured for our test topic
func initProducer() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(getKafkaAddress()),
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
		Async:    false, // Synchronous writes for testing
	}
}

// initConsumer creates a new Kafka consumer configured for our test topic
func initConsumer() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{getKafkaAddress()},
		Topic:    topicName,
		GroupID:  groupName,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}
