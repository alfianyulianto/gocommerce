package file_upload

import (
	"fmt"
	"mime/multipart"
	"slices"
)

type FileValidator struct {
	Config *Config
}

func NewFileValidator(config *Config) *FileValidator {
	return &FileValidator{Config: config}
}

func (f *FileValidator) Validate(file *multipart.FileHeader) error {
	if file.Size > f.Config.MaxSize {
		return fmt.Errorf("file size exceeds limit")
	}

	if !f.isAllowedType(file.Header.Get("Content-Type")) {
		return fmt.Errorf("file type not allowed")
	}
	return nil
}

func (f *FileValidator) isAllowedType(ext string) bool {
	return slices.Contains(f.Config.AllowedTypes, ext)

}
