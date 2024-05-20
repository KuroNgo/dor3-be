package excel

import (
	file_internal "clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

const maximumVocabulary = 5

func ReadFileForVocabulary(filename string) ([]file_internal.Vocabulary, error) {
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

	var vocabularies []file_internal.Vocabulary

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

			for i := 0; i < unitSize && currentRow < len(rows); i++ {
				row := rows[currentRow]
				if len(row) >= 8 {
					vocabulary := file_internal.Vocabulary{
						Word:          row[0],
						PartOfSpeech:  row[1],
						Pronunciation: row[2],
						Example:       row[3],
						ExplainVie:    row[4],
						ExplainEng:    row[5],
						ExampleVie:    row[6],
						ExampleEng:    row[7],
						FieldOfIT:     elementSheet,
						UnitLevel:     unitIndex,
					}
					vocabularies = append(vocabularies, vocabulary)
				}
				currentRow++
			}
			unitIndex++
		}
	}

	return vocabularies, nil
}
