package usecase

import (
	"context"
	"fmt"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/dto"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/entity"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/repository"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/search"
	"github.com/alfianyulianto/gocommerce/pkg/validation"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductUseCase interface {
	Create(ctx context.Context, request *dto.CreateProductRequest) (*dto.ProductResponse, error)
	FindById(ctx context.Context, id string) (*dto.ProductResponse, error)
}

type productUseCase struct {
	Repository repository.ProductRepository
	DB         *gorm.DB
	Log        *logrus.Logger
	ESSearch   search.ESSearch
	Config     *config.Config
}

func NewProductUseCase(repository repository.ProductRepository, DB *gorm.DB, log *logrus.Logger, ESSearch search.ESSearch, config *config.Config) ProductUseCase {
	return &productUseCase{Repository: repository, DB: DB, Log: log, ESSearch: ESSearch, Config: config}
}

func (p *productUseCase) Create(ctx context.Context, request *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	if err := validation.Validate.Struct(request); err != nil {
		return nil, err
	}

	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	product := &entity.Product{
		ID:          uuid.NewString(),
		Name:        request.Name,
		SKU:         request.SKU,
		Image:       request.Image,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	if err := p.Repository.Create(ctx, tx, product); err != nil {
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	go func() {
		err = p.ESSearch.Index(ctx, product)
		if err != nil {
			p.Log.WithField("product", product).WithError(err).Error("Failed to index product in Elasticsearch")
		}
	}()

	product.Image = fmt.Sprintf("%s:%d%s", p.Config.App.BaseURL, p.Config.App.Port, product.Image)

	return dto.ToProductResponse(product), nil
}

func (p *productUseCase) FindById(ctx context.Context, id string) (*dto.ProductResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	product := new(entity.Product)
	if err := p.Repository.FindById(ctx, tx, product, id); err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	product.Image = fmt.Sprintf("%s:%d%s", p.Config.App.BaseURL, p.Config.App.Port, product.Image)

	return dto.ToProductResponse(product), nil
}
