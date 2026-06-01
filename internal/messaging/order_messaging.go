package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alfianyulianto/gocommerce/internal/infrastructure/elasticsearch"
	"github.com/alfianyulianto/gocommerce/internal/infrastructure/kafka"
)

type OrderMessaging struct {
	Client   *elasticsearch.Client
	Consumer kafka.Consumer
}

func NewOrderMessaging(client *elasticsearch.Client, consumer kafka.Consumer) *OrderMessaging {
	return &OrderMessaging{Client: client, Consumer: consumer}
}

func (o *OrderMessaging) Start(ctx context.Context) {
	_ = o.Consumer.Consume(ctx, o.handler)
}

func (o *OrderMessaging) handler(key string, payload map[string]interface{}) error {
	switch key {
	case "order.index":
		return o.indexOrder(payload)
	default:
		return nil
	}
}

func (o *OrderMessaging) indexOrder(payload map[string]interface{}) error {
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
