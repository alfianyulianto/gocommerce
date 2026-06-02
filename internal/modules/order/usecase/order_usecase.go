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
	productRepo "github.com/alfianyulianto/gocommerce/internal/modules/product/repository"
	"github.com/alfianyulianto/gocommerce/pkg/validation"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderUseCase interface {
	Create(ctx context.Context, request *dto.CreateOrderRequest) (*dto.OrderResponse, error)
}

type orderUseCase struct {
	Repository        repository.OrderRepository
	ProductRepository productRepo.ProductRepository
	Log               *logrus.Logger
	Config            *config.Config
	Producer          kafka.Producer
}

func NewOrderUseCase(repository repository.OrderRepository, productRepository productRepo.ProductRepository, log *logrus.Logger, config *config.Config, producer kafka.Producer) OrderUseCase {
	return &orderUseCase{Repository: repository, ProductRepository: productRepository, Log: log, Config: config, Producer: producer}
}

func (o *orderUseCase) Create(ctx context.Context, request *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	if err := validation.Validate.Struct(request); err != nil {
		return nil, err
	}

	order := &entity.Order{
		ID:     uuid.New().String(),
		UserID: request.UserID,
		Status: entity.StatusPending,
		Note:   request.Note,
	}

	err := o.Repository.Transaction(ctx, func(tx *gorm.DB) error {
		var items []entity.OrderItem
		var totalAmount float32
		for _, item := range request.Items {
			product := new(entity2.Product)
			product.ID = item.ProductID
			if err := o.ProductRepository.DeductStock(ctx, tx, product, item.Quantity); err != nil {
				return err
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
		}

		order.TotalAmount = totalAmount

		if err := o.Repository.CreateWithItems(ctx, tx, order, items); err != nil {
			return err
		}

		return nil
	})

	event := dto.OrderToEvent(order)
	payload, err := json.Marshal(event)
	if err != nil {
		o.Log.WithField("order", order).WithError(err).Error("Failed to marshal order")
		return nil, err
	}
	if err = o.Producer.Publish(ctx, "order.events", "order.created", payload); err != nil {
		o.Log.WithError(err).Error("Failed to publish order")
		return nil, err
	}

	return dto.ToOrderResponse(order, o.Config), nil
}
