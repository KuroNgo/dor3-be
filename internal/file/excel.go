package file_internal

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

var (
	fieldOfIT      []string
	queueListRows0 chan string
)

func ReadFileForVocabulary(filename string) (*FileVocabulary, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		errors.New(err.Error())
		return nil, err
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

	// add list sheet into fieldOfIT property
	for sliceSheets := range queueListSheet {
		fieldOfIT = append(fieldOfIT, sliceSheets)
		rows, err := f.GetRows(sliceSheets)
		if err != nil {
			errors.New("fail to give each row")
			continue //skip the current sheet and continue next list sheet
		}

		queueListRows0 = make(chan string, len(rows[0]))

		for _, cell := range rows[0] {
			queueListRows0 <- cell
		}
		close(queueListRows0)

		for sliceRows := range queueListRows0 {
			if HaveVocabulary(sliceRows) {

			}
		}
	}

	metadata := &FileVocabulary{
		FieldOfIT: fieldOfIT,
	}

	return metadata, nil
}
