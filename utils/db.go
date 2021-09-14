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
			StoreID:        rand.Intn(4) + 1,
			StoreCode:      "eyewa",
			EntityID:       int(rand.Int31()),
			ParentEntityID: int(rand.Int31()),
		},
		Data:      data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
