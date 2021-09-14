package utils

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/eyewa/eyewa-go-lib/base"
)

func GenerateProductRecord() base.ProductModel {
	subtype := GenerateEventSubType()
	product := GenerateProduct(subtype)

	data, err := json.Marshal(product)
	logError(err)

	return base.ProductModel{
		ProductMeta: base.ProductMeta{
			ID:             uint(rand.Int()),
			StoreID:        1,
			StoreCode:      "eyewa",
			EntityID:       rand.Int(),
			ParentEntityID: rand.Int(),
		},
		Data:      data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
