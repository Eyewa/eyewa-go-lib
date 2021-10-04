package utils

import (
	"math/rand"
	"strings"
	"time"

	"github.com/eyewa/eyewa-go-lib/base"
	"github.com/eyewa/eyewa-go-lib/uuid"
	"github.com/goombaio/namegenerator"
)

func GenerateSimpleProduct() base.SimpleProduct {
	return base.SimpleProduct{
		GeneralProduct: GenerateGeneralProduct(base.SimpleProductType),
		Options:        []base.SimplesCustomOption{},
	}
}

func GenerateConfigurableProduct() base.ConfigurableProduct {
	return base.ConfigurableProduct{
		GeneralProduct: GenerateGeneralProduct(base.ConfigurableProductType),
		Variants: []base.ConfigurableSimpleProduct{
			GenerateConfigurableSimpleProduct(),
			GenerateConfigurableSimpleProduct(),
			GenerateConfigurableSimpleProduct(),
		},
	}
}

func GenerateGeneralProduct(productType base.EyewaProductType) base.GeneralProduct {
	name := GenerateName()

	image := base.ProductImage{
		URL:   "/" + strings.ReplaceAll(name, " ", "-"),
		Label: name,
	}

	description := base.ProductDescriptionHTML{
		HTML: name,
	}

	sp := "[{\"price\":\"AED 40.00\",\"id\":\"12916\",\"sku\":\"solmltp-opt-free-090\",\"name\":\"Opti-Free Pure Moist Solution 90 ml\",\"image\":\"https:\\/\\/cdn.eyewa.com\\/media\\/prescription\\/options\\/images\\/default\\/opti-free-puremoist-90ml_1.png\"}]"
	pr := "{\"total_count\":2,\"reviews\":[{\"title\":\"reviewsummary\",\"detail\":\"Comfortable \",\"nickname\":\"Ahmed mohmead\",\"date\":\"07\\/01\\/20\",\"vote\":[{\"label\":\"Rating\",\"percentage\":\"100\"}]},{\"title\":\"reviewsummary\",\"detail\":\"very comfortable, i reccommand them \",\"nickname\":\"Samy\",\"date\":\"04\\/08\\/19\",\"vote\":[{\"label\":\"Rating\",\"percentage\":\"100\"}]}]}"

	return base.GeneralProduct{
		EntityID:            rand.Int(),
		TypeID:              string(productType),
		SKU:                 uuid.NewString(),
		Name:                name,
		URLKey:              "/static-url",
		StoreID:             rand.Int(),
		ParentIDs:           &[]int{rand.Int()},
		ParentSKUs:          &[]string{uuid.NewString()},
		StoreCode:           "ae",
		AttributeSetID:      GenerateAttributeSetID(),
		MgsBrand:            GenerateBrand(),
		ContactLensSize:     rand.Int(),
		LensPackage:         ConvertStringToPointer("box"),
		StockStatus:         GenerateStockStatus(),
		Description:         description,
		ShortDescription:    description,
		SmallImage:          image,
		Rating:              rand.Intn(101),
		SolutionProduct:     ConvertStringToPointer(sp),
		ProductReviews:      ConvertStringToPointer(pr),
		MetaDescription:     name,
		MetaKeyword:         name,
		MetaTitle:           name,
		OptionLabels:        GenerateOptionLabels(),
		VirtualTryon:        ConvertIntToPointer(rand.Intn(1)),
		Categories:          GenerateCategories(),
		SpecialFromDate:     ConvertStringToPointer("N/A"),
		SpecialToDate:       ConvertStringToPointer("N/A"),
		Price:               GeneratePrice(),
		MediaGalleryEntries: []base.ProductMediaGalleryEntry{},
		Image:               image,
	}
}

func GenerateConfigurableSimpleProduct() base.ConfigurableSimpleProduct {
	name := GenerateName()

	return base.ConfigurableSimpleProduct{
		Attributes: []base.ConfigurableVariantAttribute{},
		Product: base.SimpleVariant{
			EntityID:        rand.Int(),
			TypeID:          string(base.ConfigurableProductType),
			SKU:             uuid.NewString(),
			Name:            name,
			StockStatus:     GenerateStockStatus(),
			URLKey:          "/variant-url-key",
			SpecialFromDate: ConvertStringToPointer("N/A"),
			SpecialToDate:   ConvertStringToPointer("N/A"),
			Price:           GeneratePrice(),
		},
	}
}

func GenerateName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	wordCount := rand.Intn(5) + 1

	var name string
	for i := 0; i < wordCount; i++ {
		name += nameGenerator.Generate() + " "
	}

	return strings.TrimSuffix(name, " ")
}

func GenerateBrand() string {
	var brands = []string{"30Sundays", "Artlife", "Babamio", "BlackOut",
		"Calvin Klein Jeans", "Carrera", "Charlie Max", "CHPO", "Fendi",
		"Lacoste", "Le Specs", "McQ", "MinkPink", "Mr. Wonderful",
		"MVMT", "Polaroid", "Pride", "Quay", "Ray-Ban", "ROAV",
		"Stella McCartney", "TOPFOXX", "Vogue",
	}

	return brands[rand.Intn(len(brands))]
}

func GenerateStockStatus() string {
	var status = []string{"IN_STOCK", "OUT_OF_STOCK"}

	return status[rand.Intn(len(status))]
}

func RandomType() string {
	var types = []string{"simple", "configurable", "virtual", "downloadable"}

	return types[rand.Intn(len(types))]
}

func GenerateStoreLocale() string {
	var locale = []string{"ae-ar", "ae-en"}

	return locale[rand.Intn(len(locale))]
}

func GenerateCategories() []base.ProductCategory {
	categories := map[int]string{0: "woman", 1: "man", 2: "child", 3: "sunglass", 4: "lens"}
	productCategories := make([]base.ProductCategory, 0)

	count := rand.Intn(5) + 1
	for i := 1; i < count; i++ {
		index := rand.Intn(len(categories))
		category := categories[index]
		productCategories = append(productCategories, base.ProductCategory{
			Name: category,
		})

		delete(categories, index)
	}

	return productCategories
}

func GenerateAttributeSetID() int {
	attributeSet := []int{4, 9, 10, 11, 12, 13, 14, 15, 16}

	return int(attributeSet[rand.Intn(len(attributeSet))])
}

func GeneratePrice() base.ProductPrice {
	currencies := []string{"USD", "UAE"}
	currency := currencies[rand.Intn(len(currencies))]

	amount := rand.Float64() + 50

	return base.ProductPrice{
		MaximalPrice: struct {
			Amount base.ProductPriceAmount `json:"amount"`
		}{
			Amount: base.ProductPriceAmount{
				Value:    amount + 50,
				Currency: currency,
			},
		},
		RegularPrice: struct {
			Amount base.ProductPriceAmount `json:"amount"`
		}{
			Amount: base.ProductPriceAmount{
				Value:    amount,
				Currency: currency,
			},
		},
		MinimalPrice: struct {
			Amount base.ProductPriceAmount `json:"amount"`
		}{
			Amount: base.ProductPriceAmount{
				Value:    amount - 50,
				Currency: currency,
			},
		},
	}
}

func GenerateOptionLabels() []byte {
	optionLabels := []string{
		`{"age":"Adult","bridge_size":"17 mm","frame_material":"Metal","frame_reference":"CK8043-015-52","frame_shape":"Square","frame_size":"Narrow (\u003c 131 mm)","frame_type":"Rimless","frame_width":"121 mm","lens_size":"52 mm","mgs_brand":"Calvin Klein","sku_location":"WD","temple_length":"140 mm","type":"Glasses"}`,
		`{"contact_lens_replacement":"3 months","contact_lens_size":"Pack of 2 lenses","contact_lens_type":"Color Contact Lenses","contact_lens_use":"Beauty","contact_lenses_brand":"Layala","diameter":"14.2 mm","mgs_brand":"Layala","sku_location":"EW","water_content":"38%"}`,
	}

	labels := optionLabels[rand.Intn(len(optionLabels))]

	return []byte(labels)
}
