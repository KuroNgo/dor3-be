package excel

import (
	file_internal "clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadFileForLessonManagementSystem(filename string) (file_internal.Course, []file_internal.Lesson, []file_internal.Unit, []file_internal.Vocabulary, error) {
	// Khởi tạo các kênh
	lessonCh := make(chan file_internal.Lesson)
	unitCh := make(chan file_internal.Unit)
	vocabularyCh := make(chan file_internal.Vocabulary)

	f, err := excelize.OpenFile(filename)
	if err != nil {
		return file_internal.Course{}, nil, nil, nil, err
	}

	// Lấy danh sách các sheet trong file
	sheetList := f.GetSheetList()
	if sheetList == nil {
		return file_internal.Course{}, nil, nil, nil, errors.New("empty sheet name")
	}

	// Goroutine xử lý bài học
	go func() {
		defer close(lessonCh)
		for i, sheetName := range sheetList {
			lesson := file_internal.Lesson{
				CourseID: "English for IT",
				Name:     sheetName,
				Level:    i,
			}
			lessonCh <- lesson
		}
	}()

	//Goroutine xử lý từng hàng trong sheet
	go func() {
		defer close(unitCh)
		defer close(vocabularyCh)
		for i, sheetName := range sheetList {
			rows, err := f.GetRows(sheetName)
			if err != nil {
				continue
			}
			for count, row := range rows {
				if count == 0 {
					continue
				}

				if len(row) >= 8 {
					// Tạo dữ liệu unit
					unit := file_internal.Unit{
						LessonID: sheetName,
						Name:     fmt.Sprintf("Unit %d", count+1),
						Level:    count + 1,
					}
					unitCh <- unit

					// Tạo dữ liệu vocabulary
					vocabulary := file_internal.Vocabulary{
						Word:          row[0],
						PartOfSpeech:  row[1],
						Pronunciation: row[2],
						Example:       row[3],
						ExplainVie:    row[4],
						ExplainEng:    row[5],
						ExampleVie:    row[6],
						ExampleEng:    row[7],
						FieldOfIT:     sheetList[i],
						UnitLevel:     count,
					}
					vocabularyCh <- vocabulary
				}
			}
		}
	}()

	// Nhận dữ liệu từ các kênh bài học và từ vựng
	for lesson := range lessonCh {
		lessons = append(lessons, lesson)
	}

	for vocabulary := range vocabularyCh {
		vocabularies = append(vocabularies, vocabulary)
	}

	// Nhận dữ liệu từ các kênh units và vocabularies
	for unit := range unitCh {
		units = append(units, unit)
	}
	// Trả về kết quả
	course := file_internal.Course{Name: "English for IT"}
	return course, lessons, units, vocabularies, nil
}
