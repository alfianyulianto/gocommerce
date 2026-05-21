package main

import (
	"fmt"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/internal"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure/elasticsearch"
	"github.com/alfianyulianto/gocommerce/pkg/validation"
)

func main() {
	log := infrastructure.NewLogger()
	cfg := config.Load(validation.Validate, log)
	db := infrastructure.NewDatabase(cfg, log)
	client := elasticsearch.NewClient(cfg, log)
	app := infrastructure.New(cfg, log, db)

	internal.Bootstrap(&internal.BootstrapConfig{
		App:    app,
		DB:     db,
		Logger: log,
		Client: client,
		Config: cfg,
	})

	err := app.Listen(fmt.Sprintf(":%d", cfg.App.Port))
	if err != nil {
		log.Panic("Error starting server:", err)
	}
}
