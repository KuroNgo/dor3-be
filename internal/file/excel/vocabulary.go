package excel

import (
	file_internal "clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

const (
	maximumUnitCount  = 10
	maximumVocabulary = 5
)

func ReadFileForVocabulary(filename string) ([]file_internal.Vocabulary, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	sheetList := f.GetSheetList()
	if sheetList == nil {
		return nil, errors.New("empty sheet name")
	}

	var vocabularies []file_internal.Vocabulary
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

			if unitCount <= maximumUnitCount && len(row) >= 8 {
				v := file_internal.Vocabulary{
					Word:          row[0],
					PartOfSpeech:  row[1],
					Pronunciation: row[2],
					Example:       row[3],
					ExplainVie:    row[4],
					ExplainEng:    row[5],
					ExampleVie:    row[6],
					ExampleEng:    row[7],
					FieldOfIT:     elementSheet,
					UnitID:        fmt.Sprintf("Unit%d", unitCount),
				}
				vocabularies = append(vocabularies, v)

				if len(vocabularies)%maximumVocabulary == 0 {
					unitCount++
				}
			}
		}
	}

	return vocabularies, nil
}
