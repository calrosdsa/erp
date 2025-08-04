package exporter

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelExporter interface {
	Export(sheet string, headers []interface{}, data [][]interface{}) (*bytes.Buffer, error)
}

type excelExporter struct{}

func NewExcelExporter() ExcelExporter {
	return &excelExporter{}
}

// func (e *excelExporter) Export(sheet string, headers []interface{}, data [][]interface{}) (res *bytes.Buffer, err error) {
// 	f := excelize.NewFile()
// 	defer func() {
// 		if err := f.Close(); err != nil {
// 			return
// 		}
// 	}()
// 	index, err := f.NewSheet(sheet)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//     f.SetActiveSheet(index)
// 	cell, err := excelize.CoordinatesToCellName(1, 1)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	err = f.SetSheetRow(sheet, cell, &headers)
// 	if err != nil {
// 		return
// 	}

// 	for i, row := range data {
// 		cell, err = excelize.CoordinatesToCellName(1, i+2)
// 		if err != nil {
// 			return 
// 		}
// 		err = f.SetSheetRow(sheet, cell, &row)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	res,err = f.WriteToBuffer()
	
// 	return 
// }

func (e *excelExporter) Export(sheet string, headers []interface{}, data [][]interface{}) (res *bytes.Buffer, err error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()

	// Create a new sheet
	index, err := f.NewSheet(sheet)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.SetActiveSheet(index)

	// Set the header row (A1, B1, C1, ...)
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		fmt.Println(err)
	}
	err = f.SetSheetRow(sheet, cell, &headers)
	if err != nil {
		return
	}

	// Set column width dynamically based on the longest header value
	for i := 0; i < len(headers); i++ {
		column,_ := excelize.ColumnNumberToName(i + 1) // Convert index to column name (e.g., 0 -> A, 1 -> B)
	
		maxLength := len(fmt.Sprintf("%v", headers[i])) // Get the length of the header

		// Loop through the data to find the longest string in the column
		for _, row := range data {
			cellValue := fmt.Sprintf("%v", row[i]) // Convert data cell value to string
			if len(cellValue) > maxLength {
				maxLength = len(cellValue) // Update max length if data value is longer
			}
		}

		// Set the column width based on the maximum length found
		// Excel generally considers width around 1 character = 1.1 units in width
		err = f.SetColWidth(sheet, column, column, float64(maxLength)*1.1)
		if err != nil {
			return
		}
	}

	// Set the data rows (A2, B2, C2, ...)
	for i, row := range data {
		cell, err = excelize.CoordinatesToCellName(1, i+2) // Start from row 2
		if err != nil {
			return
		}
		err = f.SetSheetRow(sheet, cell, &row)
		if err != nil {
			return
		}
	}

	// Write to buffer and return
	res, err = f.WriteToBuffer()
	return
}