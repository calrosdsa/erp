package teclumobility

type TecluMobilityPaths struct {
	Base         string
	ItemPrice string
	Order string
}

func NewTecluMobilityPath(base string) TecluMobilityPaths {
	return TecluMobilityPaths{
		Base:         base,
		ItemPrice: base + "/item-price",
		Order: base + "/order",
	}
}
