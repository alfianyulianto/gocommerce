package usecase

import (
	"context"
	"encoding/json"

	"github.com/alfianyulianto/gocommerce/config"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure/kafka"
	"github.com/alfianyulianto/gocommerce/internal/modules/order/dto"
	"github.com/alfianyulianto/gocommerce/internal/modules/order/entity"
	"github.com/alfianyulianto/gocommerce/internal/modules/order/repository"
	entity2 "github.com/alfianyulianto/gocommerce/internal/modules/product/entity"
	repository2 "github.com/alfianyulianto/gocommerce/internal/modules/product/repository"
	"github.com/alfianyulianto/gocommerce/pkg/validation"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderUseCase interface {
	Create(ctx context.Context, request *dto.CreateOrderRequest) (*dto.OrderResponse, error)
}

type orderUseCase struct {
	Repository        repository.OrderRepository
	ProductRepository repository2.ProductRepository
	DB                *gorm.DB
	Log               *logrus.Logger
	Config            *config.Config
	Producer          kafka.Producer
}

func NewOrderUseCase(repository repository.OrderRepository, productRepository repository2.ProductRepository, DB *gorm.DB, log *logrus.Logger, config *config.Config, producer kafka.Producer) OrderUseCase {
	return &orderUseCase{Repository: repository, ProductRepository: productRepository, DB: DB, Log: log, Config: config, Producer: producer}
}

func (o *orderUseCase) Create(ctx context.Context, request *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	if err := validation.Validate.Struct(request); err != nil {
		return nil, err
	}

	tx := o.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	order := &entity.Order{
		ID:     uuid.New().String(),
		UserID: request.UserID,
		Status: entity.StatusPending,
		Note:   request.Note,
	}

	var items []entity.OrderItem
	var totalAmount float32
	for _, item := range request.Items {
		product := new(entity2.Product)
		if err := o.ProductRepository.FindById(ctx, tx, product, item.ProductID); err != nil {
			return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
		}

		subtotal := float32(item.Quantity) * product.Price
		totalAmount += subtotal

		items = append(items, entity.OrderItem{
			ID:           uuid.NewString(),
			OrderID:      order.ID,
			ProductID:    product.ID,
			ProductName:  product.Name,
			ProductImage: product.Image,
			Quantity:     item.Quantity,
			UnitPrice:    product.Price,
			Subtotal:     subtotal,
		})

		if err := o.ProductRepository.DeductStock(ctx, tx, product.ID, item.Quantity); err != nil {
			return nil, fiber.NewError(fiber.StatusConflict, err.Error())
		}
	}

	order.TotalAmount = totalAmount
	order.Items = items

	if err := o.Repository.Create(ctx, tx, order); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	event := dto.OrderToEvent(order)
	payload, err := json.Marshal(event)
	if err != nil {
		o.Log.WithField("order", order).WithError(err).Error("Failed to marshal order")
		return nil, err
	}
	if err = o.Producer.Publish(ctx, "order.events", "order.index", payload); err != nil {
		o.Log.WithError(err).Error("Failed to publish order")
		return nil, err
	}

	return dto.ToOrderResponse(order, o.Config), nil
}
