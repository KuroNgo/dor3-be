package excel

import (
	file_internal "clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadFileForUnit(filename string) ([]file_internal.Unit, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Failed to close file: %v\n", err)
		}
	}()

	sheetList := f.GetSheetList()
	if sheetList == nil {
		return nil, errors.New("empty sheet name")
	}

	vocabularyCount := 0
	const maximumUnitCount = 10

	var units []file_internal.Unit
	for _, elementSheet := range sheetList {
		unitCount := 1 // Reset unitCount for each lesson

		rows, err := f.GetRows(elementSheet)
		if err != nil {
			return nil, err
		}

		for i, row := range rows {
			if i == 0 {
				continue
			}

			vocabularyCount++
			if vocabularyCount%5 == 0 {
				if len(row) >= 2 {
					u := file_internal.Unit{
						LessonID: elementSheet,
						Name:     fmt.Sprintf("Unit%d", unitCount),
					}
					units = append(units, u)
				}
				if unitCount <= maximumUnitCount {
					unitCount++
				}
				vocabularyCount = 0 // Reset vocabulary count
			}
		}
	}
	return units, nil
}
