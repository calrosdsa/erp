package helpers

import "fmt"

type TaxHelper struct{}

func NewTaxHelper() *TaxHelper {
	return &TaxHelper{}
}

func (h *TaxHelper) CalculateTotalWithTax(totalPrice int, tax float64) int {
	totalTax := float64(totalPrice) * (tax / 100)
	total := totalPrice + int(totalTax)
	fmt.Println("TOTAL PRICE WITH TAX",total,totalTax,totalPrice,tax)
	return total
}
