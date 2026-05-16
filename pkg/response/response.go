package response

import "github.com/gofiber/fiber/v3"

type Response struct {
	Message    string      `json:"message"`
	Data       any         `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItem   int64 `json:"total_item"`
	TotalPage   int   `json:"total_page"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

func OK(ctx fiber.Ctx, message string, data any, pagination *Pagination) error {
	return ctx.Status(fiber.StatusOK).JSON(Response{
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func Created(ctx fiber.Ctx, message string, data any) error {
	return ctx.Status(fiber.StatusCreated).JSON(Response{
		Message: message,
		Data:    data,
	})
}
