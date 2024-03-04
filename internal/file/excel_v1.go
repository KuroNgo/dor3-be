package file_internal

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

var (
	fieldOfIT      []string
	vocabulary     []string
	wordtype       []string
	spelling       []string
	vieMean        []string
	queueListRows0 chan string
)

func ReadFileForVocabularyV1(filename string) (*FileVocabulary, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// get all context of each sheet in spreadsheet with queue algorithm FIFO (first in first out)
	listSheet := f.GetSheetList()
	queueListSheet := make(chan string, len(listSheet))

	// Give each value of the first column into queue
	for _, elementSheet := range listSheet {
		queueListSheet <- elementSheet
	}
	close(queueListSheet)

	rowIndex := 1
	// add list sheet into fieldOfIT property
	for sliceSheets := range queueListSheet {
		//TODO add data into fieldOfIT
		fieldOfIT = append(fieldOfIT, sliceSheets)
		rows, err := f.GetRows(sliceSheets)
		if err != nil {
			errors.New("fail to give each row")
			continue //skip the current sheet and continue next list sheet
		}

		queueListRows0 = make(chan string, len(rows[0]))

		for _, elementCellInRow0 := range rows[0] {
			queueListRows0 <- elementCellInRow0
		}
		close(queueListRows0)

		for sliceRows := range queueListRows0 {
			if HaveVocabulary(sliceRows) {
				cellValue, err := f.GetCellValue(sliceSheets, sliceRows+fmt.Sprintf("%d", rowIndex))
				if err != nil {
					return nil, errors.New("fail to give value from each cell")
				}
				if cellValue == "" {
					break
				}

				vocabulary = append(vocabulary, cellValue)
				rowIndex++
			}

			if HaveWordType(sliceRows) {
				cellValue, err := f.GetCellValue(sliceSheets, sliceRows+fmt.Sprintf("%d", rowIndex))
				if err != nil {
					return nil, errors.New("fail to give value from each cell")
				}
				if cellValue == "" {
					break
				}

				wordtype = append(wordtype, cellValue)
				rowIndex++
			}

			if HaveSpelling(sliceRows) {
				cellValue, err := f.GetCellValue(sliceSheets, sliceRows+fmt.Sprintf("%d", rowIndex))
				if err != nil {
					return nil, errors.New("fail to give value from each cell")
				}
				if cellValue == "" {
					break
				}

				spelling = append(spelling, cellValue)
				rowIndex++
			}

		}
	}

	metadata := &FileVocabulary{
		Vocabulary: vocabulary,
		WordType:   wordtype,
		FieldOfIT:  fieldOfIT,
		Spelling:   spelling,
	}

	return metadata, nil
}
