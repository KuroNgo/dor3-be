package excel

import (
	file_internal "clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

var (
	courses      []file_internal.Course
	lessons      []file_internal.Lesson
	units        []file_internal.Unit
	vocabularies []file_internal.Vocabulary
)

func ReadFileForUnit(filename string) ([]file_internal.Unit, error) {
	// Open the Excel file
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Failed to close file: %v\n", err)
		}
	}()

	// Get the list of sheets in the Excel file
	sheetList := f.GetSheetList()
	if sheetList == nil {
		return nil, errors.New("empty sheet name")
	}

	var units []file_internal.Unit

	// Iterate through each sheet
	for _, elementSheet := range sheetList {
		// Get the rows in the current sheet
		rows, err := f.GetRows(elementSheet)
		if err != nil {
			continue
		}

		// Calculate the number of rows per unit, excluding the header
		totalRows := len(rows) - 1
		if totalRows <= 0 {
			continue
		}

		unitCount := (totalRows + maximumVocabulary - 1) / maximumVocabulary // Ceil division
		rowsPerUnit := totalRows / unitCount
		remainderRows := totalRows % unitCount

		unitIndex := 1
		currentRow := 1 // Start from the second row to skip the header
		for currentRow < len(rows) {
			unitSize := rowsPerUnit
			if remainderRows > 0 {
				unitSize++
				remainderRows--
			}

			unit := file_internal.Unit{
				LessonID: elementSheet,
				Name:     fmt.Sprintf("Unit %d", unitIndex),
				Level:    unitIndex,
			}
			units = append(units, unit)

			currentRow += unitSize
			unitIndex++
		}
	}

	return units, nil
}
