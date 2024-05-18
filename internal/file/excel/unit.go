package excel

import (
	file_internal "clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

var (
	courses      []file_internal.Course
	lessons      []file_internal.Lesson
	units        []file_internal.Unit
	vocabularies []file_internal.Vocabulary
)

func ReadFileForUnit(filename string) ([]file_internal.Unit, error) {
	//unitCh := make(chan file_internal.Unit)

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
		return nil, errors.New("empty sheet name")
	}

	for _, elementSheet := range sheetList {
		unitCount := 1 // Reset unitCount for each lesson

		rows, err := f.GetRows(elementSheet)
		if err != nil {
			continue
		}

		for i, row := range rows {
			if i == 0 {
				continue
			}

			if len(row) >= 2 && i%5 == 0 || i == 0 {
				unit := file_internal.Unit{
					LessonID: elementSheet,
					Name:     fmt.Sprintf("Unit %d", unitCount),
					Level:    unitCount,
				}

				units = append(units, unit)

				// Increase unit count for the current lesson
				unitCount++
			}
		}
	}

	//var units []file_internal.Unit
	//for unit := range unitCh {
	//	units = append(units, unit)
	//}

	return units, nil
}
