package order

import (
	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure/kafka"
	"github.com/alfianyulianto/gocommerce/internal/modules/order/delivery/http"
	"github.com/alfianyulianto/gocommerce/internal/modules/order/delivery/http/route"
	orderRepo "github.com/alfianyulianto/gocommerce/internal/modules/order/repository"
	"github.com/alfianyulianto/gocommerce/internal/modules/order/usecase"
	productRepo "github.com/alfianyulianto/gocommerce/internal/modules/product/repository"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Module struct {
	*http.OrderHandler
}

func NewModule(db *gorm.DB, log *logrus.Logger, config *config.Config, producer kafka.Producer) *Module {
	orderRepository := orderRepo.NewOrderRepository(db)
	productRepository := productRepo.NewProductRepository(db)
	useCase := usecase.NewOrderUseCase(orderRepository, productRepository, log, config, producer)
	handler := http.NewOrderHandler(useCase)

	return &Module{OrderHandler: handler}
}

func (m *Module) Register(router fiber.Router) {
	route.RegisterOrderRoutes(router, m.OrderHandler)
}
