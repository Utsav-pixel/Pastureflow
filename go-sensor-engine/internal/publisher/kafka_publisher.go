package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Utsav-pixel/go-sensor-engine/internal/sim"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
	batch  []kafka.Message
	mutex  sync.Mutex
}

func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
	})
	return &KafkaPublisher{
		writer: writer,
		batch:  make([]kafka.Message, 0, 100),
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

func (k *KafkaPublisher) PublishBatch(ctx context.Context, telemetries []sim.Telemetry) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	messages := make([]kafka.Message, len(telemetries))
	for i, t := range telemetries {
		value, err := json.Marshal(t)
		if err != nil {
			return err
		}
		messages[i] = kafka.Message{
			Key:   []byte(t.ZoneID),
			Value: value,
			Time:  time.Now(),
		}
	}
	return k.writer.WriteMessages(ctx, messages...)
}

func (k *KafkaPublisher) Close() error {
	fmt.Println("Closing Kafka publisher")
	return k.writer.Close()
}
