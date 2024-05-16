package excel

import (
	"clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadFileForLesson(filename string) ([]file_internal.Lesson, error) {
	lessonCh := make(chan file_internal.Lesson)

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
		return nil, errors.New("no sheets found in the Excel file")
	}

	go func() {
		defer close(lessonCh)
		for i, sheetName := range sheetList {
			lesson := file_internal.Lesson{
				CourseID: "English for IT",
				Name:     sheetName,
				Level:    i + 1,
			}
			lessonCh <- lesson
		}
	}()

	var lessons []file_internal.Lesson
	for lesson := range lessonCh {
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}
