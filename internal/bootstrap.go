package internal

import (
	"github.com/alfianyulianto/gocommerce/internal/infrastucture/elasticsearch"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	*fiber.App
	*gorm.DB
	*logrus.Logger
	*elasticsearch.Client
}

func Bootstrap(config *BootstrapConfig) {
	_ := config.App.Group("/api/v1")
}
