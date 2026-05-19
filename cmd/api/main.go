package main

import (
	"fmt"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/internal"
	"github.com/alfianyulianto/gocommerce/internal/infrastucture"
	"github.com/alfianyulianto/gocommerce/internal/infrastucture/elasticsearch"
	"github.com/alfianyulianto/gocommerce/pkg/validation"
)

func main() {
	log := infrastucture.NewLogger()
	cfg := config.Load(validation.Validate, log)
	db := infrastucture.NewDatabase(cfg, log)
	client := elasticsearch.NewClient(cfg, log)
	app := infrastucture.New(cfg, log, db)

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
