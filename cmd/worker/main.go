package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure/elasticsearch"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure/kafka"
	"github.com/alfianyulianto/gocommerce/internal/messaging"
	"github.com/alfianyulianto/gocommerce/pkg/validation"
	"github.com/sirupsen/logrus"
)

func main() {
	log := infrastructure.NewLogger()
	cfg := config.Load(validation.Validate, log)

	log.Info("Worker service started")

	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go RunOrderConsumer(cfg, log, ctx)

	stop := false
	for !stop {
		select {
		case <-quit:
			defer cancel()
			log.Info("Worker service shutting down")
			stop = true
		}
	}
}

func RunOrderConsumer(cfg *config.Config, log *logrus.Logger, ctx context.Context) {
	client := elasticsearch.NewClient(cfg, log)
	consumer := kafka.NewConsumer(cfg, log)
	orderMessaging := messaging.NewOrderMessaging(client, consumer)
	orderMessaging.Start(ctx)
}
