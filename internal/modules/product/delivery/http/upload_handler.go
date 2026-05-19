package http

import (
	"fmt"

	"github.com/alfianyulianto/gocommerce/internal/shared/file_upload"
	"github.com/alfianyulianto/gocommerce/pkg/response"
	"github.com/gofiber/fiber/v3"
)

type UploadHandler struct {
	Service *file_upload.UploadService
}

func NewUploadHandler(service *file_upload.UploadService) *UploadHandler {
	return &UploadHandler{Service: service}
}

func (u *UploadHandler) UploadImage(ctx fiber.Ctx) error {
	filename, err := u.Service.Upload(ctx, "image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	fileUrl := fmt.Sprintf("/uploads/products/%s", filename)

	return response.Created(ctx, "success upload image", fiber.Map{"filename": filename, "url": fileUrl})
}
