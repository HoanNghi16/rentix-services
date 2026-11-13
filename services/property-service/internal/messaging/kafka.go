package messaging

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisherFromEnv() (*KafkaPublisher, error) {
	brokersValue := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	if brokersValue == "" {
		brokersValue = "localhost:9092"
	}
	var brokers []string
	for _, broker := range strings.Split(brokersValue, ",") {
		broker = strings.TrimSpace(broker)
		if broker != "" {
			brokers = append(brokers, broker)
		}
	}
	if len(brokers) == 0 {
		return nil, errors.New("KAFKA_BROKERS must contain at least one broker")
	}
	return &KafkaPublisher{writer: &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Balancer:               &kafka.Hash{},
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: false,
		ReadTimeout:            10 * time.Second,
		WriteTimeout:           10 * time.Second,
	}}, nil
}

func (publisher *KafkaPublisher) Publish(ctx context.Context, topic, key string, payload []byte) error {
	return publisher.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	})
}

func (publisher *KafkaPublisher) Close() error { return publisher.writer.Close() }
