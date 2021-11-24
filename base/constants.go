package base

// List of EyewaProduct types - based on Magento's
const (
	SimpleProductType       EyewaProductType = "simple"
	ConfigurableProductType EyewaProductType = "configurable"
	DownloadableProductType EyewaProductType = "downloadable"
	VirtualProductType      EyewaProductType = "virtual"
	BundledProductType      EyewaProductType = "bundled"
	GroupedProductType      EyewaProductType = "grouped"

	ProductEnabled                ProductStatus     = 1
	ProductDisabled               ProductStatus     = 2
	ProductNotVisibleIndividually ProductVisibility = 1
	ProductVisibleCatalog         ProductVisibility = 2
	ProductVisibleSearch          ProductVisibility = 3
	ProductVisibleCatalogSearch   ProductVisibility = 4
)
