package base

import (
	"context"
	"encoding/json"
	"time"
)

// EyewaProductType
type EyewaProductType string

// EyewaEvent a base representation of an event fired/received
type EyewaEvent struct {
	ID           string `json:"id"`                                // can be used for tracing
	Name         string `json:"name"`                              // name of event - ProductUpdated, ProductDeleted etc
	EventType    string `json:"event_type"`                        // type of event's entity - Product, Order etc
	StoreLocale  string `json:"store_locale" binding:"omitempty"`  // store locale for store sa-sone, kw-ar, sa-en etc
	EventSubType string `json:"event_subtype" binding:"omitempty"` // simple, configurable", // Would be empty for category events

	// a representation on an error. provides reasons when a message ends up back
	// in the queue
	Errors []Error `json:"errors" binding:"omitempty"`

	Payload   json.RawMessage `json:"payload"`    // actual event payload
	CreatedAt string          `json:"created_at"` // time in RFC3339 format
}

// MagentoProductEvent a representation of a product event in Magento
type MagentoProductEvent struct {
	ID           string  `json:"id"`                         // uuid
	Name         string  `json:"event"`                      // name of event - product.created, catalog.created etc
	EventType    string  `json:"event_type"`                 // type of event's entity - Product, Order etc
	StoreCode    string  `json:"store_code"`                 // store code for store eyewa_kwd, eyewa_sasone etc
	StoreLocale  string  `json:"store_locale"`               // store locale for store sa-sone, kw-ar, sa-en etc
	CreatedAt    string  `json:"created_at"`                 // time in RFC3339 format of when event ocurred
	EntityID     int     `json:"entity_id"`                  // ID of the product in magento
	WebsiteID    int     `json:"website_id"`                 // ID of website store is assigned to
	StoreID      int     `json:"store_id"`                   // ID of the store the product/category belongs to
	EventSubType string  `json:"event_subtype"`              // product-simple/product-simple-custom/product-configurable", // Would be empty for category events
	Errors       []Error `json:"errors" binding:"omitempty"` // provides reasons why a message ended up in the deadletter queue for e.g
}

// Error a structural info about an error within the ecosystem
type Error struct {
	ErrorCode    int    `json:"error_code"`    // custom or http code should suffice
	ErrorMessage string `json:"error_message"` // error being reported
	CreatedAt    string `json:"created_at"`    // time in RFC3339 format
}

// EyewaEventError a structural error info about an event that failed
// during processing by a consumer - usecase: for publishing to a deadletter queue
// It is paramount to keep records of an event that failed and why.
type EyewaEventError struct {
	Event        string `json:"event"`         // string representation of event that was consumed off the queue but failed.
	ErrorMessage string `json:"error_message"` // error being reported
	CreatedAt    string `json:"created_at"`    // time in RFC3339 format
}

// MessageBrokerCallbackFunc all broker clients should define this callback fn
// so as to react to the state of events published/consumed - success/failure
type MessageBrokerCallbackFunc func(ctx context.Context, event *EyewaEvent, err error) error

// MessageBrokerMagentoProductCallbackFunc all broker clients should define this callback fn
// so as to react to the state of magento product events published/consumed - success/failure
type MessageBrokerMagentoProductCallbackFunc func(ctx context.Context, event *MagentoProductEvent) error

// EyewaProduct definition of an eyewa Product
type EyewaProduct struct {
	ID                    string
	Name                  string                  `json:"name"`                     // name of product
	SKU                   string                  `json:"sku"`                      // identifies a product
	Barcode               string                  `json:"barcode"`                  // barcode for product - multiple products can have the same barcode
	MagentoID             string                  `json:"magento_id"`               // auto-incremental id of product in Magento
	Pricing               EyewaPricing            `json:"pricing"`                  // pricing model/rules for product
	ParentID              int64                   `json:"parent_id"`                // if product is assigned to a configurable - the eyewa_product_id of the configurable. Note not MagentoID
	Brand                 string                  `json:"brand"`                    // the brand name of the product
	Store                 EyewaStore              `json:"store"`                    // the store the product is assigned to
	Visibility            []string                `json:"visibility"`               // list of where product is allowed to appear in
	Type                  string                  `json:"type"`                     // type of product - simple/configurable/virtual/downloadable etc
	Categories            []string                `json:"categories"`               // list of magento category_ids product is assigned to - listed in order of hierarchy
	ShortDescription      string                  `json:"short_description"`        // short description of product in store + store_locale
	LongDescription       string                  `json:"long_description"`         // long description of product in store + store_locale
	IsActive              bool                    `json:"is_active"`                // if produdct is active or not
	InStock               bool                    `json:"in_stock"`                 // if product is in stock or not
	IsManaged             bool                    `json:"is_managed"`               // if product is a managed product
	Rating                int64                   `json:"rating"`                   // rating of product as received from Magento - note Magento stores this in percentage
	ImageURL              string                  `json:"image_url"`                // CDN location of product image
	MagentoAttributeSetID int64                   `json:"magento_attribute_set_id"` // use to identify the type of product from Magento's perspective - Sunglass, Glasses, etc
	Attributes            []EyewaProductAttribute `json:"attributes"`               // attributes assigned to product - weight, color etc
	SEO                   ProductSEO              `json:"seo"`                      // seo info - a merge of meta and seo details in Magento
	Options               EyewaProductOptions     `json:"options"`                  // product options info from Magento
	Reviews               []ProductReview         `json:"reviews"`                  // product reviews
	MediaAssets           ProductMediaAssets      `json:"media_assets"`             // media assets for product
	CreatedAt             time.Time               // product creation datetime - Note not Magento's
	UpdatedAt             time.Time               // product updated datetime - Note not Magento's
}

// EyewaStore definition of an eyewa store
type EyewaStore struct {
	StoreCode   string `json:"store_code"`   // Magento's store_code for store
	StoreLocale string `json:"store_locale"` // Magento's store_locale for store
}

// EyewaPricing pricing details for product
type EyewaPricing struct {
	Price        string            `json:"price"`         // numerical value of price
	Currency     string            `json:"currency"`      // currency symbol of price for store
	SpecialPrice json.RawMessage   `json:"special_price"` // as receieved from Magento
	TierPrices   []json.RawMessage `json:"tier_prices"`   // as received from Magento
}

// Locale a given entity's definition for specified locale
// Should be used for cases where an entity is stored based on a locale
// or requires localization in order to have context/usefulness
type Locale struct {
	Code  string `json:"locale"` // locale code as deemed fit - en, us_en, us-en
	Value string `json:"value"`  // the value for such entity in the locale specified
}

// MagentoProductAttribute definition for product attributes in storage.
// All attributes are solely based on Magento's.
// If a product property is not deemed as primary to the definition
// of an EyewaProduct, it is regarded as an attribute and managed
// as such.
//
// {
// 		"name": "size",
// 		"magento_attribute_id": 168,
// 		"labels": [
// 			{
// 				"locale": "en",
// 				"value": "Size"
// 			},
// 			{
// 				"locale": "ar",
// 				"value": "بحجم"
// 			}
// 		]
// }
type MagentoProductAttribute struct {
	Name               string   `json:"name"`                 // name of attribute
	MagentoAttributeID int64    `json:"magento_attribute_id"` // id of attribute in Magento
	Labels             []Locale `json:"labels"`               // attribute label/name per locale
}

// EyewaProductAttribute definition of a product attribute
// Used for generating a list of attributes on product request
// for a given locale
//
// "attributes": [
// 		{
// 			"name": "Size",
// 			"magento_attribute_id": 168,
// 			"attribute_value": {
// 					"value": "Small",
// 					"magento_attribute_id": 200
// 			}
// 		},
// 		{
// 			"name": "Color",
// 			"magento_attribute_id": 177,
// 			"attribute_value": {
// 					"value": "Red",
// 					"magento_attribute_id": 213
// 			}
// 		}
// ]
type EyewaProductAttribute struct {
	Name               string                `json:"name"`                 // name of attribute
	MagentoAttributeID int64                 `json:"magento_attribute_id"` // attribute's id in Magento
	AttributeValue     MagentoAttributeValue `json:"attribute_value"`      // value of attribute
}

// MagentoAttributeValue definition for a field linked to Magento and its value
// Should be used in cases where a value to be stored requires locale info
// as well as reference to id of entity/value in Magento.
type MagentoAttributeValue struct {
	Locale             string `json:"locale"` // locale code as deemed fit - en, us_en, us-en
	Value              string `json:"value"`
	MagentoAttributeID int64  `json:"magento_attribute_id"`
}

// ProductSEO seo definition for a product
// Maps to a merge of seo and meta information in Magento
type ProductSEO struct {
	ProductURL string `json:"product_url"`
	Locales    []struct {
		Locale     string `json:"locale"`
		ProductURL string `json:"product_url"`
		MetaData   struct {
			Title       string `json:"title"`
			Keywords    string `json:"keywords"`
			Description string `json:"description"`
		} `json:"meta_data"`
	} `json:"locales"`
}

// EyewaProductOptions maps to product options in Magento
type EyewaProductOptions struct {
	Options json.RawMessage `json:"options"` // as received from Magento
}

// ProductMediaAssets maps to images and videos in Magento
type ProductMediaAssets struct {
	Assets json.RawMessage `json:"assets"`
}

// ProductReview definition of a product review
type ProductReview struct {
	Store     EyewaStore `json:"store"`      // store review was published
	Title     string     `json:"title"`      // review title
	Review    string     `json:"review"`     // review details
	Author    string     `json:"author"`     // review author
	Rating    int        `json:"rating"`     // review rating
	CreatedAt string     `json:"created_at"` // review created ts in RFC3339 format
}
