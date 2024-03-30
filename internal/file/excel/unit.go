package excel

import (
	file_internal "clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadFileForUnit(filename string) ([]file_internal.Unit, error) {
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

	var units []file_internal.Unit
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) >= 2 {
			u := file_internal.Unit{
				LessonID: row[0],
				Name:     row[1],
				Content:  row[2],
			}
			units = append(units, u)
		}
	}

	return units, nil
}
