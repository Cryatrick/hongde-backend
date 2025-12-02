package thirdparty

import(
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
	
	"hongde_backend/internal/middleware"

)

type Header struct {
	Text  string  // Header text (e.g., "ID", "Name").
	Width float64 // Column width for the header.
}

type Style struct {
	Cell      string // Specific cell (e.g., "A1", "B2")
	CellRange string // Range of cells (e.g., "A1:B2")
	FontColor string // Font color (e.g., "FF0000" for red, "00FF00" for green)
	Border    bool   // Whether to apply a full border
}

func GenerateExcelFile(excelHeader []Header, excelData []map[string]interface{}, styleArray []Style, sheetName,savePath string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			middleware.LogError(err,"Failed To Close Excel File")
			// log.Printf("Error closing Excel file: %v", err)
		}
	}()
	// Customize header cell style (e.g., bold font).
	boldStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	// Create Border Style 
	borderStyleID, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	// Create Green Style
	// greenStyleId,_ := f.NewStyle(&excelize.Style {
	// 	Font: &excelize.Font{
	// 		Color: "008000",
	// 	},
	// })


	// Create a new sheet.
	index, err := f.NewSheet(sheetName)
	if err != nil {
		middleware.LogError(err,"Failed To Create Sheet")
		return err
	}

	// Write excelHeader to the first row and set column widths.
	for colIndex, header := range excelHeader {
		colName, _ := excelize.ColumnNumberToName(colIndex + 1) // e.g., "A", "B", etc.
		cell := colName + "1"                                  // e.g., "A1", "B1", etc.

		// Write header text.
		f.SetCellValue(sheetName, cell, header.Text)

		f.SetCellStyle(sheetName, cell, cell, boldStyle)

		// Set column width for the header.
		f.SetColWidth(sheetName, colName, colName, header.Width)
	}

	// Write excelData starting from the second row.
	i := 2
	for _, row := range excelData {
		for colIndex, header := range excelHeader {
			colName, _ := excelize.ColumnNumberToName(colIndex + 1)
			cell := colName + strconv.Itoa(i) // e.g., "A2", "B2", etc.
			f.SetCellValue(sheetName, cell,  row[header.Text])
		}
		i++
	}

	// Apply custom styles from styleArray.
	// for _, style := range styleArray {

	// 	// Apply the style to the specified cell or cell range.
	// 	if style.Cell != "" {
	// 		f.SetCellStyle(sheetName, style.Cell, style.Cell, styleID)
	// 	} else if style.CellRange != "" {
	// 		f.SetCellStyle(sheetName, style.CellRange, style.CellRange, styleID)
	// 	}
	// }

	// Apply a full border to the entire table.
	lastColName, _ := excelize.ColumnNumberToName(len(excelHeader)) // Last column name
	lastRow := len(excelData) + 1                                  // Last row number
	tableRange := fmt.Sprintf("A1:%s%d", lastColName, lastRow)     // e.g., "A1:C5"

	f.SetCellStyle(sheetName, tableRange, tableRange, borderStyleID)

	// Set the active sheet.
	f.SetActiveSheet(index)

	// Delete the default "Sheet1".
	if sheetName != "Sheet1" {
		if err := f.DeleteSheet("Sheet1"); err != nil {
			// return err
			middleware.LogError(err,"Failed To Delete Default Sheet")
			return err
		}
	}

	// Save the file.
	if err := f.SaveAs(savePath); err != nil {
		middleware.LogError(err,"Failed To Save File")
		return err
	}

	return nil
}

// ReadExcelFile reads data from an Excel file and returns headers and data.
func ReadExcelFile(sheetName,filePath string) ([]Header, []map[string]interface{}, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		middleware.LogError(err,"Failed To open Excel File")
		return nil, nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			middleware.LogError(err,"Failed To Close Excel File")
			// log.Printf("Error closing Excel file: %v", err)
		}
	}()

	// Get all the rows in the first sheet.
	rows, err := f.GetRows(sheetName)
	if err != nil {
		middleware.LogError(err,"Failed to get rows from excel file")
		return nil, nil, err
	}

	// Extract headers and their widths.
	var headers []Header
	for colIndex, headerText := range rows[0] {
		colName, _ := excelize.ColumnNumberToName(colIndex + 1)
		width, _ := f.GetColWidth(sheetName, colName)
		headers = append(headers, Header{
			Text:  headerText,
			Width: width,
		})
	}

	// Extract data starting from the second row.
	var data []map[string]interface{}
	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		rowData := make(map[string]interface{})
		for colIndex, header := range headers {
			rowData[header.Text] = rows[rowIndex][colIndex]
		}
		data = append(data, rowData)
	}

	return headers, data, nil
}