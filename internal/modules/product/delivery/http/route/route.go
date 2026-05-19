package route

import (
	"github.com/alfianyulianto/gocommerce/internal/modules/product/delivery/http"
	"github.com/gofiber/fiber/v3"
)

func RegisterProductRouters(router fiber.Router, handler *http.ProductHandler, uploadHandler *http.UploadHandler) {
	group := router.Group("/products")
	group.Post("/", handler.Create)
	group.Get(":id", handler.FindById)

	group.Post("/images", uploadHandler.UploadImage)
}
