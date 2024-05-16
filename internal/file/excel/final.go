package excel

import (
	file_internal "clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadFileForLessonManagementSystem(filename string) (file_internal.Course, []file_internal.Lesson, []file_internal.Unit, []file_internal.Vocabulary, error) {
	// Khởi tạo map để lưu trữ các bài học, đơn vị và từ vựng theo tên của bảng tính
	lessonsMap := make(map[string]*file_internal.Lesson)
	unitsMap := make(map[string]map[int]*file_internal.Unit)
	vocabulariesMap := make(map[string]map[int][]file_internal.Vocabulary)

	var (
		course       = file_internal.Course{Name: "English for IT"}
		lessons      []file_internal.Lesson
		units        []file_internal.Unit
		vocabularies []file_internal.Vocabulary
	)

	f, err := excelize.OpenFile(filename)
	if err != nil {
		return course, nil, nil, nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Failed to close file: %v\n", err)
		}
	}()

	// Lấy danh sách các sheet trong file
	sheetList := f.GetSheetList()
	if sheetList == nil {
		return course, nil, nil, nil, errors.New("empty sheet name")
	}

	// Xử lý từng sheet
	for _, sheetName := range sheetList {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			continue
		}

		// Kiểm tra xem bảng tính có dữ liệu từ vựng không
		if len(rows) <= 1 {
			continue
		}

		// Tạo hoặc cập nhật thông tin cho bài học
		if _, ok := lessonsMap[sheetName]; !ok {
			lesson := file_internal.Lesson{
				CourseID: "English for IT",
				Name:     sheetName,
				Level:    len(lessons) + 1, // Sử dụng số lượng bài học hiện tại để xác định cấp độ
			}
			lessonsMap[sheetName] = &lesson
			lessons = append(lessons, lesson)
		}

		// Tạo hoặc cập nhật thông tin cho từng đơn vị và từ vựng
		unitCount := 1
		for i, row := range rows {
			if i == 0 || len(row) < 8 {
				continue
			}

			// Tạo hoặc cập nhật thông tin cho đơn vị
			if _, ok := unitsMap[sheetName]; !ok {
				unitsMap[sheetName] = make(map[int]*file_internal.Unit)
			}
			if _, ok := unitsMap[sheetName][unitCount]; !ok {
				unit := file_internal.Unit{
					LessonID: sheetName,
					Name:     fmt.Sprintf("Unit %d", unitCount),
					Level:    unitCount,
				}
				unitsMap[sheetName][unitCount] = &unit
				units = append(units, unit)
			}

			// Tạo hoặc cập nhật thông tin cho từ vựng
			if _, ok := vocabulariesMap[sheetName]; !ok {
				vocabulariesMap[sheetName] = make(map[int][]file_internal.Vocabulary)
			}
			vocab := file_internal.Vocabulary{
				Word:          row[0],
				PartOfSpeech:  row[1],
				Pronunciation: row[2],
				Example:       row[3],
				ExplainVie:    row[4],
				ExplainEng:    row[5],
				ExampleVie:    row[6],
				ExampleEng:    row[7],
				FieldOfIT:     sheetName,
				UnitLevel:     unitCount,
			}
			vocabulariesMap[sheetName][unitCount] = append(vocabulariesMap[sheetName][unitCount], vocab)

			// Nếu đạt đến maximumVocabulary, tăng số lượng đơn vị và đặt lại vocabCount
			if (i+1)%maximumVocabulary == 0 {
				unitCount++
			}
		}
	}

	// Lấy thông tin từ các map và gán vào slices tương ứng
	for _, lesson := range lessonsMap {
		lessons = append(lessons, *lesson)
	}

	for _, unitMap := range unitsMap {
		for _, unit := range unitMap {
			units = append(units, *unit)
		}
	}

	for _, vocabMap := range vocabulariesMap {
		for _, vocabSlice := range vocabMap {
			vocabularies = append(vocabularies, vocabSlice...)
		}
	}

	return course, lessons, units, vocabularies, nil
}
