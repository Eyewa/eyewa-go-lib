package utils

import (
	"math/rand"
	"time"

	"github.com/eyewa/eyewa-go-lib/base"
)

func GenerateProductRecord() base.ProductModel {
	subtype := GenerateEventSubType()
	product := GenerateProductPayload(subtype)

	return base.ProductModel{
		ProductMeta: base.ProductMeta{
			StoreID:        rand.Intn(4) + 1,
			StoreCode:      "eyewa",
			EntityID:       int(rand.Int31()),
			ParentEntityID: int(rand.Int31()),
		},
		Data:      product,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
