package file_internal

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadFileForVocabularyV2(filename string) (*FileVocabulary, error) {
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

	//index := 0

	// add list sheet into fieldOfIT property
	for sliceSheets := range queueListSheet {
		//TODO add data into fieldOfIT
		fieldOfIT = append(fieldOfIT, sliceSheets)

	}

	metadata := &FileVocabulary{
		Vocabulary: vocabulary,
		FieldOfIT:  fieldOfIT,
	}

	return metadata, nil
}
