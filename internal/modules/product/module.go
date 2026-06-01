package product

import (
	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure/elasticsearch"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/delivery/http"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/delivery/http/route"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/repository"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/search"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/usecase"
	"github.com/alfianyulianto/gocommerce/internal/shared/file_upload"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Module struct {
	*http.ProductHandler
	*http.UploadHandler
}

func NewModule(db *gorm.DB, log *logrus.Logger, client *elasticsearch.Client, config *config.Config) *Module {
	esSearch := search.NewEsProductSearch(client)
	productRepository := repository.NewProductRepository(db)
	useCase := usecase.NewProductUseCase(productRepository, log, esSearch, config)
	productHandler := http.NewProductHandler(useCase)

	fileValidator := file_upload.NewFileValidator(file_upload.ProductFileConfig())
	service := file_upload.NewUploadService(fileValidator, file_upload.ProductFileConfig())
	uploadHandler := http.NewUploadHandler(service)

	return &Module{ProductHandler: productHandler, UploadHandler: uploadHandler}
}

func (m *Module) Register(router fiber.Router) {
	route.RegisterProductRouters(router, m.ProductHandler, m.UploadHandler)
}
