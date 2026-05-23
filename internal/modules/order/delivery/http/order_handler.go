package http

import (
	"github.com/alfianyulianto/gocommerce/internal/modules/order/dto"
	"github.com/alfianyulianto/gocommerce/internal/modules/order/usecase"
	"github.com/alfianyulianto/gocommerce/pkg/response"
	"github.com/gofiber/fiber/v3"
)

type OrderHandler struct {
	UseCase usecase.OrderUseCase
}

func NewOrderHandler(useCase usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{UseCase: useCase}
}

func (o *OrderHandler) Create(ctx fiber.Ctx) error {
	request := new(dto.CreateOrderRequest)
	if err := ctx.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body")
	}

	order, err := o.UseCase.Create(ctx.Context(), request)
	if err != nil {
		return err
	}

	return response.OK(ctx, "success create order", order, nil)
}
