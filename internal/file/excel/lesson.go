package excel

import (
	"clean-architecture/internal/file"
	"errors"
	"github.com/xuri/excelize/v2"
)

func ReadFileForLesson(filename string) ([]file_internal.Lesson, error) {
	// khởi tạo channel
	lessonCh := make(chan file_internal.Lesson)

	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	sheetList := f.GetSheetList()
	if sheetList == nil {
		return nil, errors.New("no sheets found in the Excel file")
	}

	go func() {
		defer close(lessonCh)
		for i, elementSheet := range sheetList {
			l := file_internal.Lesson{
				CourseID: "English for IT",
				Name:     elementSheet,
				Level:    i,
			}
			lessonCh <- l
		}
	}()

	// nhận dữ liệu từ các kênh bai học
	for lesson := range lessonCh {
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}
