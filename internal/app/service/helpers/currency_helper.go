package helpers

import (
	"erp/api/common"
	"fmt"
	"math"
)

type CurrencyHelper interface {
	FloatToInt(value float64) int64
	IntToFloat(value int64) float64
	FormatCurrencyAmount(amount int64,currency string) string
	CurrencyExchange(amount int64,exchangeRate int32) int64
	SplitFloat(num float64) (int, int)
	
}

type currencyHelper struct {
}

func NewCurrencyHelper() CurrencyHelper {
	return &currencyHelper{}
}
func (h *currencyHelper)SplitFloat(num float64) (int, int) {
	// Get the absolute value of the input number
	absValue := math.Abs(num)

	// Extract the integer part
	integerPart := int(absValue)

	// Extract the decimal part
	decimalPart := int(math.Round((absValue - float64(integerPart)) * 100))

	return integerPart, decimalPart
}
func (h *currencyHelper)FormatCurrencyAmount(amount int64,currency string) string {
	v := h.IntToFloat(amount)
	symbol := h.getCurrencySymbol(currency)
	return fmt.Sprintf("%s %.2f",symbol,v)
}
func (h *currencyHelper) CurrencyExchange(amount int64,exchangeRate int32) int64 {
	fmt.Println("CurrencyExchange",amount,exchangeRate)
	if exchangeRate == 1 {
		return amount
	}
	// fmt.Println("RATE",rate)
	return (amount * int64(exchangeRate)) / 100
}
//Convert From 99.99 to 9999
func (h *currencyHelper) FloatToInt(value float64) int64 {
	return int64(math.Round(value * 100))
}

//Convert From 9999 to 99.99
func (h *currencyHelper) IntToFloat(value int64) float64 {
	floatValue := float64(value) / 100
	return floatValue
}


func (h *currencyHelper) getCurrencySymbol(currency string) string{
	switch currency {
	case string(common.CurrencyCodeUSD):
		return "$"
	case string(common.CurrencyCodeBOB):
		return "Bs"
	default:
		return currency
	}
}