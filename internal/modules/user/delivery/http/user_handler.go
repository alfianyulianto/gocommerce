package http

import (
	"github.com/alfianyulianto/gocommerce/internal/modules/user/dto"
	"github.com/alfianyulianto/gocommerce/internal/modules/user/usecase"
	"github.com/alfianyulianto/gocommerce/pkg/response"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	UseCase usecase.UserUseCase
}

func NewUserHandler(useCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{UseCase: useCase}
}

func (h *UserHandler) Create(ctx fiber.Ctx) error {
	request := new(dto.CreateUserRequest)
	if err := ctx.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body")
	}

	user, err := h.UseCase.Create(ctx.Context(), request)
	if err != nil {
		return err
	}

	return response.Created(ctx, "success created user", user)
}

func (h *UserHandler) Update(ctx fiber.Ctx) error {
	request := new(dto.UpdateUserRequest)
	if err := ctx.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body")
	}

	request.ID = ctx.Params("id")

	user, err := h.UseCase.Update(ctx.Context(), request)
	if err != nil {
		return err
	}

	return response.OK(ctx, "success updated user", user, nil)
}

func (h *UserHandler) FindById(ctx fiber.Ctx) error {
	user, err := h.UseCase.FindById(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}

	return response.OK(ctx, "success find user", user, nil)
}

func (h *UserHandler) Delete(ctx fiber.Ctx) error {
	err := h.UseCase.Delete(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}

	return response.OK(ctx, "success deleted user", nil, nil)
}

func (h *UserHandler) FindAll(ctx fiber.Ctx) error {
	request := new(dto.UserFilterRequest)

	request.Page = fiber.Query(ctx, "page", 1)
	request.PerPage = fiber.Query(ctx, "per_page", 10)
	request.OrderBy = ctx.Query("order_by", "created_at")
	request.OrderDirection = ctx.Query("order_direction", "desc")

	users, pagination, err := h.UseCase.FindAll(ctx.Context(), request)
	if err != nil {
		return err
	}

	return response.OK(ctx, "success find users", users, pagination)
}

func (h *UserHandler) Search(ctx fiber.Ctx) error {
	request := new(dto.SearchUserRequest)

	request.Search = ctx.Query("q")
	request.Page = fiber.Query(ctx, "page", 1)
	request.PerPage = fiber.Query(ctx, "per_page", 10)
	request.OrderBy = ctx.Query("order_by", "created_at")
	request.OrderDirection = ctx.Query("order_direction", "desc")

	users, pagination, err := h.UseCase.Search(ctx.Context(), request)
	if err != nil {
		return err
	}

	return response.OK(ctx, "success search users", users, pagination)
}
