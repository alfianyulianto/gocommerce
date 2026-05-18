package route

import (
	"github.com/alfianyulianto/gocommerce/internal/modules/user/delivery/http"
	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(router fiber.Router, handler *http.UserHandler) {
	group := router.Group("/users")
	group.Post("/", handler.Create)
	group.Patch(":id", handler.Update)
	group.Get("/search", handler.Search)
	group.Get(":id", handler.FindById)
	group.Delete(":id", handler.Delete)
	group.Get("/", handler.FindAll)
}
