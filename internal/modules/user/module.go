package user

import (
	"github.com/alfianyulianto/gocommerce/internal/infrastucture/elasticsearch"
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
	esSearch := search.NewEsUserSearch(client)
	userRepository := repository.NewUserRepository()
	useCase := usecase.NewUserUseCase(userRepository, db, log, esSearch)
	userHandler := http.NewUserHandler(useCase)
	return &Module{UserHandler: userHandler}
}

func (m *Module) Register(router fiber.Router) {
	route.RegisterUserRoutes(router, m.UserHandler)
}
