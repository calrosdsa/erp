package currency

import (
	// "fmt"
	"math"
	"strings"
)

var unitsMap = map[int]string{
	1: "un",
	2: "dos",
	3: "tres",
	4: "cuatro",
	5: "cinco",
	6: "seis",
	7: "siete",
	8: "ocho",
	9: "nueve",
}

var tensMap = map[int]string{
	2: "veinti",
	3: "treinta",
	4: "cuarenta",
	5: "cincuenta",
	6: "sesenta",
	7: "setenta",
	8: "ochenta",
	9: "noventa",
}

var specials = map[int]string{
	10: "diez",
	11: "once",
	12: "doce",
	13: "trece",
	14: "catorce",
	15: "quince",
	20: "veinte",
	100: "cien",
}

func convertLessThanThousand(n int) string {
	if n == 0 {
		return ""
	}
	
	var result string
	
	// Centenas
	if hundreds := n / 100; hundreds > 0 {
		switch hundreds {
		case 1:
			if n == 100 {
				return "cien"
			}
			result += "ciento "
		case 5:
			result += "quinientos "
		case 7:
			result += "setecientos "
		case 9:
			result += "novecientos "
		default:
			result += unitsMap[hundreds] + "cientos "
		}
		n %= 100
	}
	
	// Decenas y unidades
	if n == 0 {
		return strings.TrimSpace(result)
	}
	
	if val, ok := specials[n]; ok {
		result += val
		return result
	}
	
	if n >= 16 && n <= 19 {
		return result + "dieci" + unitsMap[n-10]
	}
	
	if n >= 21 && n <= 29 {
		return result + "veinti" + unitsMap[n-20]
	}
	
	tens := n / 10 * 10
	units := n % 10
	
	if tens > 0 {
		result += tensMap[tens/10]
	}
	
	if units > 0 {
		if tens >= 30 {
			result += " y "
		}
		result += unitsMap[units]
	}
	
	return strings.ReplaceAll(result, "  ", " ")
}

func convertNumberToWords(n int) string {
	if n == 0 {
		return "cero"
	}
	
	groups := []int{}
	for n > 0 {
		groups = append(groups, n%1000)
		n /= 1000
	}
	
	result := ""
	// scales := []string{"", "mil", "millón", "millones"}
	
	for i := len(groups) - 1; i >= 0; i-- {
		group := groups[i]
		if group == 0 {
			continue
		}
		
		groupStr := convertLessThanThousand(group)
		scale := ""
		
		switch i {
		case 1:
			scale = "mil"
			if group == 1 {
				groupStr = ""
			}
		case 2:
			if group == 1 {
				scale = "millón"
				groupStr = "un"
			} else {
				scale = "millones"
			}
		}
		
		result += groupStr + " " + scale + " "
	}
	
	result = strings.TrimSpace(result)
	result = strings.ReplaceAll(result, "  ", " ")
	
	// Correcciones gramaticales
	result = strings.ReplaceAll(result, "un mil", "mil")
	
	return result
}

func AmountToWords(amount float64) string {
	if amount == 0 {
		return "cero"
	}
	
	isNegative := amount < 0
	absolute := math.Abs(amount)
	
	integerPart := int(math.Trunc(absolute))
	decimalPart := int(math.Round((absolute - math.Trunc(absolute)) * 100))
	
	intWords := convertNumberToWords(integerPart)
	decWords := convertNumberToWords(decimalPart)
	
	result := ""
	if isNegative {
		result += "menos "
	}
	result += intWords
	
	if decimalPart > 0 {
		result += " con " + decWords + " centavos"
	}
	
	return result
}