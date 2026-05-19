package http

import (
	"github.com/alfianyulianto/gocommerce/internal/modules/product/dto"
	"github.com/alfianyulianto/gocommerce/internal/modules/product/usecase"
	"github.com/alfianyulianto/gocommerce/pkg/response"
	"github.com/gofiber/fiber/v3"
)

type ProductHandler struct {
	UseCase usecase.ProductUseCase
}

func NewProductHandler(useCase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{UseCase: useCase}
}

func (p *ProductHandler) Create(ctx fiber.Ctx) error {
	request := new(dto.CreateProductRequest)
	if err := ctx.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body")
	}

	product, err := p.UseCase.Create(ctx.Context(), request)
	if err != nil {
		return err
	}

	return response.Created(ctx, "success create user", product)
}

func (p *ProductHandler) FindById(ctx fiber.Ctx) error {
	product, err := p.UseCase.FindById(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}
	return response.OK(ctx, "success find product", product, nil)
}
