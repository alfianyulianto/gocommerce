package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alfianyulianto/gocommerce/internal/infrastructure/elasticsearch"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure/kafka"
)

type OrderConsumer struct {
	Client   *elasticsearch.Client
	Consumer kafka.Consumer
}

func NewOrderMessaging(client *elasticsearch.Client, consumer kafka.Consumer) *OrderConsumer {
	return &OrderConsumer{Client: client, Consumer: consumer}
}

func (o *OrderConsumer) Start(ctx context.Context) {
	_ = o.Consumer.Consume(ctx, o.handler)
}

func (o *OrderConsumer) handler(key string, payload map[string]interface{}) error {
	switch key {
	case "order.created":
		return o.indexOrder(payload)
	default:
		return nil
	}
}

func (o *OrderConsumer) indexOrder(payload map[string]interface{}) error {
	rawID, ok := payload["id"]
	if !ok {
		return fmt.Errorf("Missing order id in paylod")
	}

	var orderIDStr string
	switch v := rawID.(type) {
	case string:
		orderIDStr = v
	default:
		bytes, _ := json.Marshal(v)
		orderIDStr = string(bytes)
	}

	return o.Client.Index(context.Background(), o.Client.OrderIndex(), orderIDStr, payload)
}
