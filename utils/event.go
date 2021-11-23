package utils

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/eyewa/eyewa-go-lib/base"
	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/eyewa/eyewa-go-lib/uuid"
	"go.uber.org/zap"
)

const (
	ProductEventType                string = "Product"
	ProductCreated                  string = "product.created"
	ProductUpdated                  string = "product.updated"
	ProductDeleted                  string = "product.deleted"
	SimpleProductEventSubType       string = "product-" + string(base.SimpleProductType)
	ConfigurableProductEventSubType string = "product-" + string(base.ConfigurableProductType)
)

func GenerateRandomProductEvent() *base.EyewaEvent {
	eventSubType := GenerateEventSubType()

	product := GenerateProductPayload(eventSubType)

	// sample event
	return &base.EyewaEvent{
		ID:           uuid.NewString(),
		Name:         GenerateEventName(),
		EventType:    ProductEventType,
		EventSubType: string(eventSubType),
		StoreCode:    GenerateStoreCode(),
		Payload:      product,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}
}

func GenerateConfigurableProductEvent() *base.EyewaEvent {
	product := GenerateProductPayload(base.ConfigurableProductType)

	return &base.EyewaEvent{
		ID:           uuid.NewString(),
		Name:         GenerateEventName(),
		EventType:    ProductEventType,
		StoreCode:    GenerateStoreCode(),
		EventSubType: ConfigurableProductEventSubType,
		Payload:      product,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}
}

func GenerateSimpleProductEvent() *base.EyewaEvent {
	product := GenerateProductPayload(base.SimpleProductType)

	return &base.EyewaEvent{
		ID:           uuid.NewString(),
		Name:         GenerateEventName(),
		EventType:    ProductEventType,
		StoreCode:    GenerateStoreCode(),
		EventSubType: SimpleProductEventSubType,
		Payload:      product,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}
}

func GenerateEventName() string {
	names := []string{ProductCreated, ProductUpdated}
	return names[rand.Intn(len(names))]
}

func GenerateEventSubType() base.EyewaProductType {
	subTypes := []base.EyewaProductType{
		base.EyewaProductType(SimpleProductEventSubType),
		base.EyewaProductType(ConfigurableProductEventSubType),
	}
	return subTypes[rand.Intn(len(subTypes))]
}

func GenerateProductPayload(eventSubType base.EyewaProductType) []byte {
	switch eventSubType {
	case base.SimpleProductType:
		simpleProduct := GenerateSimpleProduct()

		data, err := json.Marshal(simpleProduct)
		logError(err)

		return data
	case base.ConfigurableProductType:
		configurableProduct := GenerateConfigurableProduct()

		data, err := json.Marshal(configurableProduct)
		logError(err)

		return data
	}

	return nil
}

func logError(err error) {
	if err != nil {
		log.Error("Got error", zap.Error(err))
	}
}
