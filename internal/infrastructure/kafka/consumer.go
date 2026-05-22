package kafka

import (
	"context"
	"encoding/json"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type MessageHandler func(key string, payload map[string]interface{}) error

type Consumer interface {
	Consume(ctx context.Context, handler MessageHandler) error
	Close() error
}

type consumer struct {
	Reader *kafka.Reader
	Log    *logrus.Logger
}

func NewConsumer(cfg *config.Config, log *logrus.Logger) Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupId,
		Topic:   cfg.Kafka.Topic,
	})

	return &consumer{Reader: r, Log: log}
}

func (c *consumer) Consume(ctx context.Context, handler MessageHandler) error {
	c.Log.Info("Kafka consume started")
	for {
		msg, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			c.Log.WithFields(logrus.Fields{
				"topic": msg.Topic,
				"key":   string(msg.Key),
				"value": string(msg.Value),
			}).WithError(err).Error("Kafka reader error")
			continue
		}

		if ctx.Err() != nil {
			c.Log.WithError(ctx.Err()).Info("Kafka consume stopped")
			return nil
		}

		var payload map[string]interface{}
		if err = json.Unmarshal(msg.Value, &payload); err != nil {
			c.Log.WithFields(logrus.Fields{
				"topic": msg.Topic,
				"key":   string(msg.Key),
				"value": string(msg.Value),
			}).WithError(err).Error("Failed to unmarshal message")
			continue
		}

		if err = handler(string(msg.Key), payload); err != nil {
			continue
		}

		if err = c.Reader.CommitMessages(ctx, msg); err != nil {
			c.Log.WithFields(logrus.Fields{
				"topic": msg.Topic,
				"key":   string(msg.Key),
				"value": string(msg.Value),
			}).WithError(err).Error("Kafka commit error")
		}
	}
}

func (c *consumer) Close() error {
	return c.Reader.Close()
}
