package file

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadFileExcel(filename string) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return
	}

	defer func() {
		//Close the spreadsheet
		if err := f.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	// Get all the rows
	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, row := range rows {
			for _, colCell := range row {
				fmt.Println(colCell, "\t")
			}
			fmt.Println()
		}
	}
}
