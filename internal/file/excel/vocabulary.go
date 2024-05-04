package excel

import (
	file_internal "clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

const maximumVocabulary = 5

func ReadFileForVocabulary(filename string) ([]file_internal.Vocabulary, error) {
	vocabularyCh := make(chan file_internal.Vocabulary)

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

	go func() {
		defer close(vocabularyCh)
		for i, elementSheet := range sheetList {
			unitCount := 1 // Reset unitCount for each lesson

			rows, err := f.GetRows(elementSheet)
			if err != nil {
				return
			}

			vocabCount := 0 // Initialize vocabulary count for the current unit

			for j, row := range rows {
				if j == 0 {
					continue
				}

				if len(row) >= 8 {
					v := file_internal.Vocabulary{
						Word:          row[0],
						PartOfSpeech:  row[1],
						Pronunciation: row[2],
						Example:       row[3],
						ExplainVie:    row[4],
						ExplainEng:    row[5],
						ExampleVie:    row[6],
						ExampleEng:    row[7],
						FieldOfIT:     sheetList[i],
						UnitLevel:     unitCount,
					}

					vocabularyCh <- v

					// Increase vocabulary count for the current unit
					vocabCount++

					// If vocabCount reaches the maximum vocabulary, create a new unit
					if vocabCount == maximumVocabulary {
						unitCount++
						vocabCount = 1 // Reset vocabCount for the new unit
					}
				}
			}
		}
	}()

	for vocabulary := range vocabularyCh {
		vocabularies = append(vocabularies, vocabulary)
	}

	return vocabularies, nil
}
