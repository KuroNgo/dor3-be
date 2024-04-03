package excel

import (
	file_internal "clean-architecture/internal/file"
	"errors"
	"github.com/xuri/excelize/v2"
)

func ReadFileForMean(filename string) ([]file_internal.Mean, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	sheetList := f.GetSheetList()
	if sheetList == nil {
		return nil, errors.New("empty sheet name")
	}

	var means []file_internal.Mean

	for _, elementSheet := range sheetList {
		rows, err := f.GetRows(elementSheet)
		if err != nil {
			return nil, err
		}

		for i, row := range rows {
			if i == 0 {
				continue
			}

			if len(row) >= 8 {
				m := file_internal.Mean{
					LessonID:     elementSheet,
					VocabularyID: row[0],
					ExplainVie:   row[4],
					ExplainEng:   row[5],
					ExampleVie:   row[6],
					ExampleEng:   row[7],
				}
				means = append(means, m)
			}
		}
	}

	return means, nil
}
