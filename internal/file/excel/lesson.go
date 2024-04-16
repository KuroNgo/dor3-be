package excel

import (
	"clean-architecture/internal/file"
	"errors"
	"github.com/xuri/excelize/v2"
	"sync"
)

func ReadFileForLesson(filename string) ([]file_internal.Lesson, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	sheetList := f.GetSheetList()
	if sheetList == nil {
		return nil, errors.New("no sheets found in the Excel file")
	}

	var lessons []file_internal.Lesson
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, elementSheet := range sheetList {
		wg.Add(1)
		go func(sheetName string, level int) {
			defer wg.Done()

			l := file_internal.Lesson{
				CourseID: "English for IT",
				Name:     sheetName,
				Level:    level,
			}

			mu.Lock()
			lessons = append(lessons, l)
			mu.Unlock()
		}(elementSheet, i)
	}

	wg.Wait()

	return lessons, nil
}
