package excel

import (
	"clean-architecture/internal/file"
	"errors"
	"github.com/xuri/excelize/v2"
)

func ReadFileForLesson(filename string) ([]file_internal.Lesson, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	sheetList := f.GetSheetList()
	if sheetList == nil {
		return nil, errors.New("empty sheet name")
	}

	var lessons []file_internal.Lesson
	for i, elementSheet := range sheetList {
		l := file_internal.Lesson{
			CourseID: "English for IT",
			Name:     elementSheet,
			Content:  "null",
			Level:    i,
		}

		lessons = append(lessons, l)
	}

	return lessons, nil
}
