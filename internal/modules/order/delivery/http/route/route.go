package route

import (
	"github.com/alfianyulianto/gocommerce/internal/modules/order/delivery/http"
	"github.com/gofiber/fiber/v3"
)

func RegisterOrderRoutes(router fiber.Router, handler *http.OrderHandler) {
	group := router.Group("/orders")
	group.Post("/", handler.Create)
}
