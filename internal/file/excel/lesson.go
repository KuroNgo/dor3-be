package excel

import (
	"clean-architecture/internal/file"
	"errors"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func ReadFileForLesson(filename string) ([]file_internal.Lesson, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, errors.New("empty sheet name")
	}

	var lessons []file_internal.Lesson
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) >= 2 {
			level, err := strconv.Atoi(row[3])
			if err != nil {
				continue
			}
			l := file_internal.Lesson{
				CourseID: row[0],
				Name:     row[1],
				Content:  row[2],
				Level:    level,
			}
			lessons = append(lessons, l)
		}
	}

	return lessons, nil
}
