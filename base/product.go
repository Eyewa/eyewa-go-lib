package base

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// ProductModel product information saved to storage
//
// This data can either be for a configurable or a simple product
// It is not important at the stage of saving a product to storage in
// maintaining its structure as seen in ConfigurableProduct or
// SimpleProduct.
//
// The MVP objective here is to save 'a product' irrespective of type,
// separate fields such as StoreID, ParentSKU etc, as respective db columns
// for quick query lookups, and have the rest of the product data as a JSON blob.
type ProductModel struct {
	ProductMeta

	// The data contained here is a typical Magento Product marshalled as a
	// JSON blob and conforms to Magento's GraphQL expected data response.
	// i.e either ConfigurableProduct or SimpleProduct as a JSON blob
	Data datatypes.JSON `json:"data"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeleteAt  time.Time
}

// ProductMeta these are fields internal to the service either
// there for lookup or assists during a transformation process
type ProductMeta struct {
	ID              uint   `gorm:"primaryKey" json:"-"`
	StoreID         int    `gorm:"index:uix_pdt_store_entity, unique"`
	StoreCode       string `gorm:"index:uix_pdt_store_entity, unique"`
	EntityID        int    `gorm:"index:uix_pdt_store_entity, unique"`
	TypeID          string
	ParentEntityIDs *datatypes.JSON // if simple, list of configurables it is assigned to
}

// ConfigurableProduct magento's configurable product definition
type ConfigurableProduct struct {
	GeneralProduct
	Variants       []ConfigurableSimpleProduct `json:"variants"`
	JSONConfigData JSONConfigData              `json:"jsonConfigData"`
}

// SimpleProduct magento's simple product definition
type SimpleProduct struct {
	GeneralProduct
	Options []SimplesCustomOption `json:"options"`
}

// ConfigurableSimpleProduct definition for simples embedded in a configurable
// used when a configurable is requested as a collection of variants
//
// {
//     "id": 13697,
//     "type_id": "configurable",
//     "variants": [
//         {
//             "attributes": [
//                 {
//                     "code": "sphere",
//                     "label": "0.00",
//                     "value_index": 2725
//                 },
//                 {
//                     "code": "contact_lens_color_name",
//                     "label": "Cloudy Gray",
//                     "value_index": 2591
//                 }
//             ],
//             "product": {
//                 "id": 13699,
//                 "type_id": "simple",
//                 "sku": "cost3mo-lay-lens-p02-cloudy-gray-sp0000",
//                 "name": "Layala Lenses - 2 Lenses - Cloudy Gray",
//                 "stock_status": "IN_STOCK",
//                 "mgs_brand": 2598,
//                 "url_key": "layala-lenses-pack-of-2-cloudy-gray",
//                 "virtual_tryon": null,
//                 "special_from_date": null,
//                 "special_to_date": null,
//                 "price": {
//                     "maximalPrice": {
//                         "amount": {
//                             "currency": "AED",
//                             "value": 149
//                         }
//                     },
//                     "regularPrice": {
//                         "amount": {
//                             "currency": "AED",
//                             "value": 149
//                         }
//                     },
//                     "minimalPrice": {
//                         "amount": {
//                             "currency": "AED",
//                             "value": 149
//                         }
//                     }
//                 }
//             }
//         }
//     ]
// }
type ConfigurableSimpleProduct struct {
	Attributes []ConfigurableVariantAttribute `json:"attributes"`
	Product    SimpleVariant                  `json:"product"`
}

// ConfigurableVariantAttribute attributes data for configurable variants
type ConfigurableVariantAttribute struct {
	Code       string `json:"code"`
	Label      string `json:"label"`
	ValueIndex int    `json:"value_index"`
}

// SimpleVariant definition for a variant under a configurable data
type SimpleVariant struct {
	EntityID            int                        `json:"id"`
	TypeID              string                     `json:"type_id"`
	SKU                 string                     `json:"sku"`
	Name                string                     `json:"name"`
	StockStatus         string                     `json:"stock_status"`
	URLKey              string                     `json:"url_key"`
	SpecialPrice        *float64                   `json:"special_price"`
	SpecialFromDate     *string                    `json:"special_from_date"`
	SpecialToDate       *string                    `json:"special_to_date"`
	Price               ProductPrice               `json:"price"`
	MediaGalleryEntries []ProductMediaGalleryEntry `json:"media_gallery_entries"`
}

// GeneralProduct a typical definition of a product common to both configurables or simples.
type GeneralProduct struct {
	EntityID            int                        `json:"id"`
	TypeID              string                     `json:"type_id"`
	SKU                 string                     `json:"sku"`
	Name                string                     `json:"name"`
	URLKey              string                     `json:"url_key"`
	StoreID             int                        `json:"store_id"`
	ParentIDs           *[]int                     `json:"parent_ids"`
	ParentSKUs          *[]string                  `json:"parent_skus"`
	StoreCode           string                     `json:"store_code"`
	AttributeSetID      int                        `json:"attribute_set_id"`
	AttributeSetName    string                     `json:"attribute_set_name"`
	MgsBrand            string                     `json:"mgs_brand"`
	ContactLensSize     int                        `json:"contact_lens_size"`
	LensPackage         *string                    `json:"lens_package"`
	StockStatus         string                     `json:"stock_status"`
	Description         ProductDescriptionHTML     `json:"description"`
	ShortDescription    ProductDescriptionHTML     `json:"short_description"`
	SmallImage          ProductImage               `json:"small_image"`
	ThumbnailImage      ProductImage               `json:"thumbnail_image"`
	Rating              int                        `json:"rating"`
	SolutionProduct     *string                    `json:"solution_product"` // stringified json.RawMessage
	ProductReviews      *string                    `json:"productReviews"`   // stringified ProductReviews
	MetaDescription     string                     `json:"meta_description"`
	MetaKeyword         string                     `json:"meta_keyword"`
	MetaTitle           string                     `json:"meta_title"`
	OptionLabels        *string                    `json:"option_labels"` // stringified OptionLabels
	VirtualTryon        *int                       `json:"virtual_tryon"`
	Categories          []ProductCategory          `json:"categories"`
	SpecialPrice        *float64                   `json:"special_price"`
	SpecialFromDate     *string                    `json:"special_from_date"`
	SpecialToDate       *string                    `json:"special_to_date"`
	Price               ProductPrice               `json:"price"`
	MediaGalleryEntries []ProductMediaGalleryEntry `json:"media_gallery_entries"`
	Image               ProductImage               `json:"image"`
}

// ProductCategory product category definition
type ProductCategory struct {
	Name     string `json:"name"`
	ID       int    `json:"id"`
	URLKey   string `json:"url_key"`
	Position int    `json:"position"`
}

// ProductReview product review info
type ProductReview struct {
	Title    string        `json:"title"`
	Detail   string        `json:"detail"`
	Nickname string        `json:"nickname"`
	Date     string        `json:"date"`
	Votes    []ProductVote `json:"vote"`
}

// ProductReviews a collection of reviews submitted for a product
//
// {
//     "id": 13697,
//     "type_id": "configurable",
//     "attribute_set_id": 10,
//     "sku": "cost3mo-lay-lens-p02",
//     "product_reviews": [
//         {
//             "title": "Amazing color",
//             "detail": "Belle Elite Silky Gold looks amaaazing.  Beautiful color. Really lovely!!! ❤️❤️",
//             "nickname": "Nouf",
//             "date": "17/01/18",
//             "vote": [
//                 {
//                     "label": "Rating",
//                     "percentage": "100"
//                 }
//             ]
//         },
//         {
//             "title": "Amazing color",
//             "detail": "Belle Elite Silky Gold looks amaaazing.  Beautiful color. Really lovely!!! ❤️❤️",
//             "nickname": "Nouf",
//             "date": "17/01/18",
//             "vote": [
//                 {
//                     "label": "Rating",
//                     "percentage": "100"
//                 }
//             ]
//         }
//     ]
// }
type ProductReviews struct {
	TotalCount int             `json:"total_count"`
	Reviews    []ProductReview `json:"reviews"`
}

// ProductVote product vote definition
type ProductVote struct {
	Label      string `json:"label"`
	Percentage int    `json:"percentage"`
}

// ProductPrice product pricing definition
type ProductPrice struct {
	MaximalPrice struct {
		Amount ProductPriceAmount `json:"amount"`
	} `json:"maximalPrice"`
	RegularPrice struct {
		Amount ProductPriceAmount `json:"amount"`
	} `json:"regularPrice"`
	MinimalPrice struct {
		Amount ProductPriceAmount `json:"amount"`
	} `json:"minimalPrice"`
}

// ProductPriceAmount price amount definition
type ProductPriceAmount struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

// ProductMediaGalleryEntry product media gallery definition
type ProductMediaGalleryEntry struct {
	ID           int           `json:"id"`
	Label        *string       `json:"label"`
	Position     int           `json:"position"`
	File         string        `json:"file"`
	Disabled     bool          `json:"disabled"`
	MediaType    string        `json:"media_type"`
	VideoContent *ProductVideo `json:"video_content"`
}

// ProductImage a product image definition
type ProductImage struct {
	URL   string `json:"url"`
	Label string `json:"label"`
}

// ProductDescriptionHTML product html description
type ProductDescriptionHTML struct {
	HTML string `json:"html"`
}

// SimplesCustomOption product options for simples
type SimplesCustomOption struct {
	OptionID  int                        `json:"option_id"`
	Required  bool                       `json:"required"`
	SortOrder int                        `json:"sort_order"`
	Title     string                     `json:"title"`
	Value     []SimplesCustomOptionValue `json:"value,omitempty"`
}

// SimplesCustomOptionValue product option values for product options
type SimplesCustomOptionValue struct {
	Price        float64 `json:"price"`
	PriceType    string  `json:"price_type"`
	SKU          *string `json:"sku"`
	OptionTypeID int     `json:"option_type_id"`
	Title        string  `json:"title"`
	SortOrder    int     `json:"sort_order"`
}

// ProductVideo a product's video content
type ProductVideo struct {
	MediaType   *string `json:"media_type,omitempty"`
	Description *string `json:"video_description,omitempty"`
	MetaData    *string `json:"video_metadata,omitempty"`
	Provider    *string `json:"video_provider,omitempty"`
	Title       *string `json:"video_title,omitempty"`
	URL         *string `json:"video_url,omitempty"`
}

// SolutionProduct solution product for a product
// This definition gets marshalled into GeneralProduct.SolutionProduct
// as a json.RawMessage.
type SolutionProduct struct {
	Price string `json:"price"` // price with currency symbol
	ID    int    `json:"id"`
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type JSONConfigData struct {
	ChildProducts   json.RawMessage `json:"ChildProducts"`
	SuperAttributes json.RawMessage `json:"SuperAttributes"`
}

// TableName overrides the table name for ProductModel
func (ProductModel) TableName() string {
	return "products"
}
