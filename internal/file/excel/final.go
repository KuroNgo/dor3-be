package excel

import (
	file_internal "clean-architecture/internal/file"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadFileForLessonManagementSystem(filename string) ([]file_internal.Final, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	sheetList := f.GetSheetList()
	if sheetList == nil {
		return nil, errors.New("empty sheet name")
	}

	const maximumUnitCount = 10
	var final []file_internal.Final

	for _, elementSheet := range sheetList {
		unitCount := 1 // Reset unitCount for each lesson

		rows, err := f.GetRows(elementSheet)
		if err != nil {
			return nil, err
		}

		for i, row := range rows {
			if i == 0 {
				continue
			}

			// Handle Vocabulary
			if len(row) >= 4 {
				v := file_internal.Final{
					LessonCourseID:          "English for IT",
					LessonName:              elementSheet,
					LessonLevel:             i,
					LessonContent:           "",
					VocabularyWord:          row[0],
					VocabularyPartOfSpeech:  row[1],
					VocabularyPronunciation: row[2],
					VocabularyExample:       row[3],
					VocabularyFieldOfIT:     elementSheet,
					VocabularyUnitID:        fmt.Sprintf("Unit%d", unitCount),
				}
				final = append(final, v)

				if len(final)%maximumUnitCount == 0 {
					unitCount++
				}
			}

			// Handle Mean
			if len(row) >= 8 {
				m := file_internal.Final{
					MeanLessonID:     "English for IT",
					LessonName:       elementSheet,
					MeanVocabularyID: row[0],
					MeanExplainVie:   row[4],
					MeanExplainEng:   row[5],
					MeanExampleVie:   row[6],
					MeanExampleEng:   row[7],
				}
				final = append(final, m)
			}
		}
	}

	return final, nil
}
