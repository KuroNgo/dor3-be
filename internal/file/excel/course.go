package excel

import (
	"clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadFileForCourse(filename string) ([]file_internal.Course, error) {
	f, err := excelize.OpenFile(filename)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil { // Kiểm tra lỗi khi đóng tệp
			fmt.Printf("Failed to close file: %v\n", err)
		}
	}()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, errors.New("empty sheet name")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) >= 2 {
			c := file_internal.Course{
				Name:        row[0],
				Description: row[1],
			}
			courses = append(courses, c)
		}
	}

	return courses, nil
}
