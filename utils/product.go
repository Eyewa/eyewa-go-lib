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
	}
}

func GenerateConfigurableProduct() base.ConfigurableProduct {
	variants := []base.ConfigurableSimpleProduct{
		GenerateConfigurableSimpleProduct(),
		GenerateConfigurableSimpleProduct(),
		GenerateConfigurableSimpleProduct(),
	}
	return base.ConfigurableProduct{
		GeneralProduct: GenerateGeneralProduct(base.ConfigurableProductType),
		Variants:       variants,
		JSONConfigData: GenerateJSONConfigData(variants),
	}
}

func GenerateJSONConfigData(variants []base.ConfigurableSimpleProduct) base.JSONConfigData {
	childProducts := new(string)
	*childProducts = `[{
		"contact_lens_size": "422",
		"contact_lens_size_value": "Pack of 12 lenses",
		"entity_id": "35628",
		"image": "https://cdn.eyewa.com/600x600/media/catalog/product/placeholder/default/eyewa-placeholder-555.jpg",
		"name": "أكيوفيو أوسايس هيدرا كلير بلس - عبوة من 12",
		"price": "315.0000",
		"qty": "0.0000",
		"sku": "` + variants[0].Product.SKU + `",
		"small_image": "https://cdn.eyewa.com/195x195/media/catalog/product/placeholder/default/eyewa-placeholder-555.jpg",
		"special_price": "265.0000",
		"sphere": "2700",
		"sphere_value": "-3.75",
		"stock_status": "1",
		"thumbnail": "https://cdn.eyewa.com/80x80/media/catalog/product/placeholder/default/eyewa-placeholder-555.jpg",
		"visibility": "1"
	  },
	  {
		"contact_lens_size": "422",
		"contact_lens_size_value": "Pack of 12 lenses",
		"entity_id": "35630",
		"image": "https://cdn.eyewa.com/600x600/media/catalog/product/placeholder/default/eyewa-placeholder-555.jpg",
		"name": "أكيوفيو أوسايس هيدرا كلير بلس - عبوة من 12",
		"price": "315.0000",
		"qty": "0.0000",
		"sku": "` + variants[1].Product.SKU + `",
		"small_image": "https://cdn.eyewa.com/195x195/media/catalog/product/placeholder/default/eyewa-placeholder-555.jpg",
		"special_price": "265.0000",
		"sphere": "2702",
		"sphere_value": "-4.25",
		"stock_status": "1",
		"thumbnail": "https://cdn.eyewa.com/80x80/media/catalog/product/placeholder/default/eyewa-placeholder-555.jpg",
		"visibility": "4"
	  }]`

	superAttributes := new(string)
	*superAttributes = `{
		"198": {
		  "id": 198,
		  "code": "contact_lens_size",
		  "swatch": null,
		  "label": "حجم العبوة"
		},
		"214": {
		  "id": 214,
		  "code": "sphere",
		  "swatch": null,
		  "label": "درجة قصر النظر"
		}
	  }`

	return base.JSONConfigData{
		ChildProducts:   childProducts,
		SuperAttributes: superAttributes,
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
	pr := "{\"total_count\":2,\"reviews\":[{\"title\":\"reviewsummary\",\"detail\":\"Comfortable \",\"nickname\":\"Ahmed mohmead\",\"date\":\"07\\/01\\/20\",\"vote\":[{\"label\":\"Rating\",\"percentage\":100}]},{\"title\":\"reviewsummary\",\"detail\":\"very comfortable, i reccommand them \",\"nickname\":\"Samy\",\"date\":\"04\\/08\\/19\",\"vote\":[{\"label\":\"Rating\",\"percentage\":100}]}]}"

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
		ThumbnailImage:      image,
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
	brands := []string{
		"30Sundays", "Artlife", "Babamio", "BlackOut",
		"Calvin Klein Jeans", "Carrera", "Charlie Max", "CHPO", "Fendi",
		"Lacoste", "Le Specs", "McQ", "MinkPink", "Mr. Wonderful",
		"MVMT", "Polaroid", "Pride", "Quay", "Ray-Ban", "ROAV",
		"Stella McCartney", "TOPFOXX", "Vogue",
	}

	return brands[rand.Intn(len(brands))]
}

func GenerateStockStatus() string {
	status := []string{"IN_STOCK", "OUT_OF_STOCK"}

	return status[rand.Intn(len(status))]
}

func RandomType() string {
	types := []string{"simple", "configurable", "virtual", "downloadable"}

	return types[rand.Intn(len(types))]
}

func GenerateStoreCode() string {
	locale := []string{"ae-ar", "ae-en"}

	return locale[rand.Intn(len(locale))]
}

func GenerateCategories() []base.ProductCategory {
	categories := []string{"woman", "man", "child", "sunglass", "lens"}
	productCategories := make([]base.ProductCategory, 0)

	for index, category := range categories {
		productCategories = append(productCategories, base.ProductCategory{
			Name:     category,
			ID:       index,
			URLKey:   category,
			Position: index,
		})
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

func GenerateOptionLabels() *string {
	optionLabels := []string{
		"{\"age\":{\"id\":8329,\"label\":\"للكبار\",\"attribute_label\":\"Age\"},\"bridge_size\":{\"id\":1080,\"label\":\"19 mm\",\"attribute_label\":\"Bridge Size\"},\"frame_color\":{\"id\":[2073],\"label\":[\"بني\"],\"attribute_label\":\"Frame Color\",\"hex_code\": [\"#0099cc\"]},\"frame_finish\":{\"id\":8292,\"label\":\"لامع\",\"attribute_label\":\"Frame Finish\"},\"frame_material\":{\"id\":339,\"label\":\"بلاستيك\",\"attribute_label\":\"Frame Material\"},\"frame_reference\":{\"id\":564,\"label\":\"CK8559-236-54\",\"attribute_label\":\"Frame Reference\"},\"frame_shape\":{\"id\":344,\"label\":\"مربع\",\"attribute_label\":\"Frame Shape\"},\"frame_size\":{\"id\":348,\"label\":\"صغير (\\u003c 131 مم)\",\"attribute_label\":\"Frame Size\"},\"frame_type\":{\"id\":351,\"label\":\"إطار كامل\",\"attribute_label\":\"Frame Type\"},\"frame_width\":{\"id\":1102,\"label\":\"127 mm\",\"attribute_label\":\"Frame Width\"},\"gender\":{\"id\":355,\"label\":\"نساء\",\"attribute_label\":\"Gender\"},\"lens_size\":{\"id\":1055,\"label\":\"54 mm\",\"attribute_label\":\"Lens Size\"},\"mgs_brand\":{\"id\":7,\"label\":\"كالفن كلاين\",\"attribute_label\":\"Brand\"},\"sku_location\":{\"id\":1043,\"label\":\"WD\",\"attribute_label\":\"sku_location\"},\"temple_length\":{\"id\":1094,\"label\":\"140 mm\",\"attribute_label\":\"Temple Length\"},\"type\":{\"id\":357,\"label\":\"نظارات طبية\",\"attribute_label\":\"Type\"}}",
		"{\"age\":{\"id\":8329,\"label\":\"للكبار\",\"attribute_label\":\"Age\"},\"bridge_size\":{\"id\":1078,\"label\":\"17 mm\",\"attribute_label\":\"Bridge Size\"},\"frame_color\":{\"id\":[2072],\"label\":[\"أزرق\"],\"attribute_label\":\"Frame Color\",\"hex_code\": [\"#0099cc\"]},\"frame_material\":{\"id\":338,\"label\":\"معدن\",\"attribute_label\":\"Frame Material\"},\"frame_reference\":{\"id\":715,\"label\":\"L2223-424-56\",\"attribute_label\":\"Frame Reference\"},\"frame_shape\":{\"id\":344,\"label\":\"مربع\",\"attribute_label\":\"Frame Shape\"},\"frame_size\":{\"id\":348,\"label\":\"صغير (\\u003c 131 مم)\",\"attribute_label\":\"Frame Size\"},\"frame_type\":{\"id\":351,\"label\":\"إطار كامل\",\"attribute_label\":\"Frame Type\"},\"frame_width\":{\"id\":1106,\"label\":\"131 mm\",\"attribute_label\":\"Frame Width\"},\"gender\":{\"id\":355,\"label\":\"نساء\",\"attribute_label\":\"Gender\"},\"lens_size\":{\"id\":1057,\"label\":\"56 mm\",\"attribute_label\":\"Lens Size\"},\"mgs_brand\":{\"id\":10,\"label\":\"لاكوست\",\"attribute_label\":\"Brand\"},\"sku_location\":{\"id\":1043,\"label\":\"WD\",\"attribute_label\":\"sku_location\"},\"temple_length\":{\"id\":1097,\"label\":\"145 mm\",\"attribute_label\":\"Temple Length\"},\"type\":{\"id\":357,\"label\":\"نظارات طبية\",\"attribute_label\":\"Type\"}}",
	}

	labels := optionLabels[rand.Intn(len(optionLabels))]

	return &labels
}
