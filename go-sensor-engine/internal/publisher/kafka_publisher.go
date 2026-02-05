package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Utsav-pixel/go-sensor-engine/internal/sim"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.Hash{},
	})
	return &KafkaPublisher{
		writer: writer,
	}
}

func (k *KafkaPublisher) Publish(ctx context.Context, t sim.Telemetry) error {
	value, err := json.Marshal(t)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(t.ZoneID),
		Value: value,
		Time:  time.Now(),
	}
	return k.writer.WriteMessages(ctx, msg)
}

func (k *KafkaPublisher) Close() error {
	fmt.Println("Closing Kafka publisher")
	return k.writer.Close()
}
