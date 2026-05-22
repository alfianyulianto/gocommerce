package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/internal"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure/elasticsearch"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure/kafka"
	"github.com/alfianyulianto/gocommerce/pkg/validation"
)

func main() {
	log := infrastructure.NewLogger()
	cfg := config.Load(validation.Validate, log)
	db := infrastructure.NewDatabase(cfg, log)
	client := elasticsearch.NewClient(cfg, log)
	app := infrastructure.New(cfg, log, db)
	producer := kafka.NewProducer(cfg, log)
	defer producer.Close()

	internal.Bootstrap(&internal.BootstrapConfig{
		App:      app,
		DB:       db,
		Logger:   log,
		Client:   client,
		Config:   cfg,
		Producer: producer,
	})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			log.WithError(err).Error("Failed to shutdown app")
		}
		log.Info("Gracefully shutting down")
	}()

	err := app.Listen(fmt.Sprintf(":%d", cfg.App.Port))
	if err != nil {
		log.Panic("Error starting server:", err)
	}
}
