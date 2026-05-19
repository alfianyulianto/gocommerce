package file_upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UploadService struct {
	Validator *FileValidator
	Config    *Config
}

func NewUploadService(validator *FileValidator, config *Config) *UploadService {
	return &UploadService{Validator: validator, Config: config}
}

func (u *UploadService) Upload(ctx fiber.Ctx, formName string) (string, error) {
	file, err := ctx.FormFile(formName)
	if err != nil {
		return "", fmt.Errorf("form file error: %v", err)
	}

	if err = u.Validator.Validate(file); err != nil {
		return "", err
	}

	if err = os.MkdirAll(u.Config.UploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	ext := filepath.Ext(file.Filename)
	safeName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(u.Config.UploadDir, safeName)

	if err = ctx.SaveFile(file, savePath); err != nil {
		return "", fmt.Errorf("save file error: %v", err)
	}

	return safeName, nil
}
