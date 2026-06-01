package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/dto"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/entity"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/repository"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/search"
	"github.com/alfianyulianto/gocommerce/pkg/response"
	"github.com/alfianyulianto/gocommerce/pkg/validation"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProductUseCase interface {
	Create(ctx context.Context, request *dto.CreateProductRequest) (*dto.ProductResponse, error)
	Update(ctx context.Context, request *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	FindById(ctx context.Context, id string) (*dto.ProductResponse, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, filter *dto.ProductFilterRequest) ([]dto.ProductResponse, *response.Pagination, error)
	Search(ctx context.Context, filter *dto.SearchProductRequest) ([]dto.ProductResponse, *response.Pagination, error)
}

type productUseCase struct {
	Repository    repository.ProductRepository
	Log           *logrus.Logger
	ProductSearch search.ProductSearch
	Config        *config.Config
}

func NewProductUseCase(repository repository.ProductRepository, log *logrus.Logger, productSearch search.ProductSearch, config *config.Config) ProductUseCase {
	return &productUseCase{Repository: repository, Log: log, ProductSearch: productSearch, Config: config}
}

func (p *productUseCase) Create(ctx context.Context, request *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	if err := validation.Validate.Struct(request); err != nil {
		return nil, err
	}

	product := &entity.Product{
		ID:          uuid.NewString(),
		Name:        request.Name,
		SKU:         request.SKU,
		Image:       request.Image,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	if err := p.Repository.Create(ctx, product); err != nil {
		return nil, err
	}

	go func() {
		err := p.ProductSearch.Index(ctx, product)
		if err != nil {
			p.Log.WithField("product", product).WithError(err).Error("Failed to index product in Elasticsearch")
		}
	}()

	product.Image = fmt.Sprintf("%s:%d%s", p.Config.App.BaseURL, p.Config.App.Port, product.Image)

	return dto.ToProductResponse(product), nil
}

func (p *productUseCase) Update(ctx context.Context, request *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	if err := validation.Validate.Struct(request); err != nil {
		return nil, err
	}

	product := new(entity.Product)
	if err := p.Repository.FindById(ctx, product, request.ID); err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	product.Name = request.Name
	product.SKU = request.SKU
	product.Description = request.Description
	product.Price = request.Price
	product.Stock = request.Stock

	if request.Image != "" {
		product.Image = request.Image
	}

	if err := p.Repository.Update(ctx, product); err != nil {
		return nil, err
	}

	go func() {
		if err := p.ProductSearch.Index(ctx, product); err != nil {
			p.Log.WithField("product", product).WithError(err).Error("Failed to index update product in Elasticsearch")
		}
	}()

	return dto.ToProductResponse(product), nil
}

func (p *productUseCase) FindById(ctx context.Context, id string) (*dto.ProductResponse, error) {
	product := new(entity.Product)
	if err := p.Repository.FindById(ctx, product, id); err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	product.Image = fmt.Sprintf("%s:%d%s", p.Config.App.BaseURL, p.Config.App.Port, product.Image)

	return dto.ToProductResponse(product), nil
}

func (p *productUseCase) Delete(ctx context.Context, id string) error {
	product := new(entity.Product)
	if err := p.Repository.FindById(ctx, product, id); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	if err := p.Repository.Delete(ctx, product); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	go func() {
		if err := p.ProductSearch.Delete(ctx, id); err != nil {
			p.Log.WithField("product", product).WithError(err).Error("Failed to delete product in Elasticsearch")
		}
	}()

	return nil
}

func (p *productUseCase) FindAll(ctx context.Context, request *dto.ProductFilterRequest) ([]dto.ProductResponse, *response.Pagination, error) {
	if err := validation.Validate.Struct(request); err != nil {
		return nil, nil, err
	}

	filter := dto.ToProductFilter(request)
	products, total, err := p.Repository.FindAll(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		productResponse := dto.ToProductResponse(&product)
		responses = append(responses, *productResponse)
	}

	totalPage := int(math.Ceil(float64(total) / float64(request.PerPage)))

	pagination := &response.Pagination{
		CurrentPage: request.Page,
		PerPage:     request.PerPage,
		TotalItem:   total,
		TotalPage:   totalPage,
		HasNext:     request.Page < totalPage,
		HasPrev:     request.Page > 1 && request.Page <= totalPage,
	}

	return responses, pagination, nil
}

func (p *productUseCase) Search(ctx context.Context, request *dto.SearchProductRequest) ([]dto.ProductResponse, *response.Pagination, error) {

	if err := validation.Validate.Struct(request); err != nil {
		return nil, nil, err
	}

	products, total, err := p.ProductSearch.Search(ctx, request)
	if err != nil {
		return nil, nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(request.PerPage)))

	pagination := &response.Pagination{
		CurrentPage: request.Page,
		PerPage:     request.PerPage,
		TotalItem:   total,
		TotalPage:   totalPage,
		HasNext:     request.Page < totalPage,
		HasPrev:     request.Page > 1 && request.Page <= totalPage,
	}
	return products, pagination, nil

}
