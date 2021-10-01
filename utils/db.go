package utils

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/eyewa/eyewa-go-lib/base"
	"gorm.io/datatypes"
)

func GenerateProductRecord() base.ProductModel {
	subtype := GenerateEventSubType()
	product := GenerateProductPayload(subtype)

	j, _ := json.Marshal(rand.Int31())
	r := datatypes.JSON(j)

	return base.ProductModel{
		ProductMeta: base.ProductMeta{
			StoreID:         rand.Intn(4) + 1,
			StoreCode:       "eyewa",
			EntityID:        int(rand.Int31()),
			ParentEntityIDs: &r,
		},
		Data:      product,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
