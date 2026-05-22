package kafka

import (
	"context"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Producer interface {
	Publish(ctx context.Context, topic, key string, payload []byte) error
	Close() error
}

type producer struct {
	Writer *kafka.Writer
	Log    *logrus.Logger
}

func NewProducer(cfg *config.Config, log *logrus.Logger) Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Kafka.Brokers...),
		Balancer:               &kafka.RoundRobin{},
		RequiredAcks:           kafka.RequireOne,
		MaxAttempts:            5,
		AllowAutoTopicCreation: true,
	}
	return &producer{Writer: w, Log: log}
}

func (p *producer) Publish(ctx context.Context, topic, key string, payload []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	}

	if err := p.Writer.WriteMessages(ctx, msg); err != nil {
		p.Log.WithFields(logrus.Fields{
			"topic": topic,
			"key":   key,
			"value": string(payload),
		}).WithError(err).Error("Failed to write messages")
		return err
	}

	p.Log.WithFields(logrus.Fields{
		"topic": topic,
		"key":   key,
		"value": string(payload),
	}).Info("Successfully published message")
	return nil
}

func (p *producer) Close() error {
	err := p.Writer.Close()
	if err != nil {
		p.Log.WithError(err).Error("Failed to close writer")
		return err
	}
	p.Log.Info("Successfully closed writer")
	return nil
}
