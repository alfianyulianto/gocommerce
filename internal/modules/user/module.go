package user

import (
	"github.com/alfianyulianto/gocommerce/internal/infrastructure/elasticsearch"
	"github.com/alfianyulianto/gocommerce/internal/modules/user/delivery/http"
	"github.com/alfianyulianto/gocommerce/internal/modules/user/delivery/http/route"
	"github.com/alfianyulianto/gocommerce/internal/modules/user/repository"
	"github.com/alfianyulianto/gocommerce/internal/modules/user/search"
	"github.com/alfianyulianto/gocommerce/internal/modules/user/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Module struct {
	*http.UserHandler
}

func NewModule(db *gorm.DB, log *logrus.Logger, client *elasticsearch.Client) *Module {
	userSearch := search.NewUserSearch(client)
	userRepository := repository.NewUserRepository(db)
	useCase := usecase.NewUserUseCase(userRepository, log, userSearch)
	userHandler := http.NewUserHandler(useCase)
	return &Module{UserHandler: userHandler}
}

func (m *Module) Register(router fiber.Router) {
	route.RegisterUserRoutes(router, m.UserHandler)
}
