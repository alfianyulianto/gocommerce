package internal

import (
	"github.com/alfianyulianto/gocommerce/internal/infrastucture/elasticsearch"
	"github.com/alfianyulianto/gocommerce/internal/modules/user"
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
	v1 := config.App.Group("/api/v1")

	userModule := user.NewModule(config.DB, config.Logger, config.Client)
	userModule.Register(v1)
}
