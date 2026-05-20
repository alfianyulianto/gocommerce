package route

import (
	"github.com/alfianyulianto/gocommerce/internal/modules/product/delivery/http"
	"github.com/gofiber/fiber/v3"
)

func RegisterProductRouters(router fiber.Router, handler *http.ProductHandler, uploadHandler *http.UploadHandler) {
	group := router.Group("/products")
	group.Post("/", handler.Create)
	group.Patch(":id", handler.Update)
	group.Get("/search", handler.Search)
	group.Get(":id", handler.FindById)
	group.Delete(":id", handler.Delete)
	group.Get("/", handler.FindAll)

	group.Post("/images", uploadHandler.UploadImage)
}
