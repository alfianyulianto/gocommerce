package infrastructure

import (
	"errors"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func New(cfg *config.Config, log *logrus.Logger, db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: cfg.App.Name,
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			return errorHandler(ctx, err, log)
		},
	})

	return app
}

func errorHandler(ctx fiber.Ctx, err error, log *logrus.Logger) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		log.WithFields(logrus.Fields{
			"path":   ctx.Path(),
			"method": ctx.Method(),
			"error":  e.Error(),
		}).Warn(e.Message)
		return ctx.Status(e.Code).JSON(fiber.Map{
			"message": e.Message,
		})
	}

	if errors.As(err, new(validator.ValidationErrors)) {
		log.WithFields(logrus.Fields{
			"path":   ctx.Path(),
			"method": ctx.Method(),
			"errors": validation.ParseMessage(err),
		}).Warn("Validation Error")
		return ctx.Status(400).JSON(fiber.Map{
			"message": "validation error",
			"errors":  validation.ParseMessage(err),
		})
	}

	log.WithError(err).Error("Internal Server Error")
	return ctx.Status(code).JSONP(fiber.Map{
		"message": "internal server error",
	})
}
