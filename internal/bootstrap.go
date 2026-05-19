package internal

import (
	"os"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/internal/infrastucture/elasticsearch"
	"github.com/alfianyulianto/gocommerce/internal/modules/product"
	"github.com/alfianyulianto/gocommerce/internal/modules/user"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	*fiber.App
	*gorm.DB
	*logrus.Logger
	*elasticsearch.Client
	*config.Config
}

func Bootstrap(config *BootstrapConfig) {
	config.App.Use("/uploads", static.New("", static.Config{
		FS:     os.DirFS("./uploads"),
		Browse: true,
	}))

	v1 := config.App.Group("/api/v1")

	userModule := user.NewModule(config.DB, config.Logger, config.Client)
	userModule.Register(v1)

	productModule := product.NewModule(config.DB, config.Logger, config.Client, config.Config)
	productModule.Register(v1)
}
