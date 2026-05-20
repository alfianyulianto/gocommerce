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

func (p *ProductHandler) Update(ctx fiber.Ctx) error {
	request := new(dto.UpdateProductRequest)
	if err := ctx.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body")
	}

	request.ID = ctx.Params("id")

	product, err := p.UseCase.Update(ctx.Context(), request)
	if err != nil {
		return err
	}

	return response.OK(ctx, "success update user", product, nil)
}

func (p *ProductHandler) FindById(ctx fiber.Ctx) error {
	product, err := p.UseCase.FindById(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}
	return response.OK(ctx, "success find product", product, nil)
}

func (p *ProductHandler) Delete(ctx fiber.Ctx) error {
	if err := p.UseCase.Delete(ctx.Context(), ctx.Params("id")); err != nil {
		return err
	}

	return response.OK(ctx, "success delete product", nil, nil)
}

func (p *ProductHandler) FindAll(ctx fiber.Ctx) error {
	request := new(dto.ProductFilterRequest)

	request.MinPrice = fiber.Query[int32](ctx, "min_price", 0)
	request.MaxPrice = fiber.Query[int32](ctx, "max_price", 0)
	request.Page = fiber.Query(ctx, "page", 1)
	request.PerPage = fiber.Query(ctx, "per_page", 10)
	request.OrderBy = ctx.Query("order_by", "created_at")
	request.OrderDirection = ctx.Query("order_direction", "desc")

	products, pagination, err := p.UseCase.FindAll(ctx.Context(), request)
	if err != nil {
		return err
	}

	return response.OK(ctx, "success find products", products, pagination)
}

func (p *ProductHandler) Search(ctx fiber.Ctx) error {
	request := new(dto.SearchProductRequest)

	request.Search = ctx.Query("q")
	request.MinPrice = fiber.Query[int32](ctx, "min_price", 0)
	request.MaxPrice = fiber.Query[int32](ctx, "max_price", 0)
	request.Page = fiber.Query(ctx, "page", 1)
	request.PerPage = fiber.Query(ctx, "per_page", 10)
	request.OrderBy = ctx.Query("order_by", "created_at")
	request.OrderDirection = ctx.Query("order_direction", "desc")

	products, pagination, err := p.UseCase.Search(ctx.Context(), request)
	if err != nil {
		return err
	}

	return response.OK(ctx, "success search products", products, pagination)
}
