package course_controller

import (
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/internal"
	file_internal "clean-architecture/internal/file"
	"clean-architecture/internal/file/excel"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"sync"
	"time"
)

func (c *CourseController) CreateOneCourse(ctx *gin.Context) {
	var courseInput course_domain.Input
	if err := ctx.ShouldBindJSON(&courseInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if err := internal.IsValidCourse(courseInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	course := &course_domain.Course{
		Id:          primitive.NewObjectID(),
		Name:        courseInput.Name,
		Description: courseInput.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := c.CourseUseCase.CreateOne(ctx, course)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (c *CourseController) CreateCourseWithFile(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(8 << 20) // 8MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("files")
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !file_internal.IsExcel(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not an Excel file",
		})
		return
	}

	err = ctx.SaveUploadedFile(file, file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		err := os.Remove(file.Filename)
		if err != nil {
			fmt.Printf("Failed to delete temporary file: %v\n", err)
		}
	}()

	result, err := excel.ReadFileForCourse(file.Filename)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var courses []course_domain.Course
	for _, course := range result {
		c := course_domain.Course{
			Id:          primitive.NewObjectID(),
			Description: course.Description,
			Name:        course.Name,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			//WhoUpdated:
		}
		courses = append(courses, c)
	}

	var wg sync.WaitGroup
	successCount := 0
	errChan := make(chan error)

	for _, course := range courses {
		wg.Add(1)
		go func(course course_domain.Course) {
			defer wg.Done()
			err := c.CourseUseCase.CreateOne(ctx, &course)
			if err != nil {
				errChan <- err
				return
			}
			successCount++
		}(course)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()

	if successCount == 0 {
		errChan <- errors.New("failed to create any course")
		return
	}

	select {
	case err := <-errChan:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error creating course: %v\n": err,
		})
	default:
		ctx.JSON(http.StatusOK, gin.H{
			"status":        "success",
			"success_count": successCount,
		})
	}

}

func (c *CourseController) CreateLessonManagementWithFile(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(8 << 20) // 8MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("files")
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !file_internal.IsExcel(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not an Excel file",
		})
		return
	}

	err = ctx.SaveUploadedFile(file, file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		err := os.Remove(file.Filename)
		if err != nil {
			fmt.Printf("Failed to delete temporary file: %v\n", err)
		}
	}()

	result, err := excel.ReadFileForLessonManagementSystem(file.Filename)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var lessons []lesson_domain.Lesson
	for _, lesson := range result {
		courseID, err := c.LessonUseCase.FindCourseIDByCourseName(ctx, lesson.LessonCourseID)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		l := lesson_domain.Lesson{
			ID:          primitive.NewObjectID(),
			CourseID:    courseID,
			Name:        lesson.LessonName,
			Level:       lesson.LessonLevel,
			IsCompleted: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			//WhoUpdates:
		}
		lessons = append(lessons, l)
	}

	for _, lesson := range lessons {
		err = c.LessonUseCase.CreateOneByNameCourse(ctx, &lesson)
		if err != nil {
			continue
		}
	}

	var units []unit_domain.Unit
	for _, unit := range result {
		lessonID, err := c.UnitUseCase.FindLessonIDByLessonName(ctx, unit.UnitLessonID)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		elUnit := unit_domain.Unit{
			ID:         primitive.NewObjectID(),
			LessonID:   lessonID,
			Name:       unit.UnitName,
			ImageURL:   "",
			IsComplete: 0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			//WhoUpdates: user.FullName,
		}
		units = append(units, elUnit)
	}

	for _, unit := range units {
		err = c.UnitUseCase.CreateOneByNameLesson(ctx, &unit)
		if err != nil {
			continue
		}
	}

	var vocabularies []vocabulary_domain.Vocabulary

	for _, vocabulary := range result {
		unitID, err := c.VocabularyUseCase.FindUnitIDByUnitLevel(ctx, vocabulary.VocabularyUnitLevel)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		v := vocabulary_domain.Vocabulary{
			Id:            primitive.NewObjectID(),
			UnitID:        unitID,
			Word:          vocabulary.VocabularyWord,
			PartOfSpeech:  vocabulary.VocabularyPartOfSpeech,
			Pronunciation: vocabulary.VocabularyPronunciation,
			ExplainEng:    vocabulary.MeanExplainEng,
			ExplainVie:    vocabulary.MeanExplainVie,
			ExampleVie:    vocabulary.MeanExampleVie,
			ExampleEng:    vocabulary.MeanExampleEng,
			FieldOfIT:     vocabulary.VocabularyFieldOfIT,
			LinkURL:       "",
		}
		vocabularies = append(vocabularies, v)
	}

	for _, vocabulary := range vocabularies {
		err = c.VocabularyUseCase.CreateOneByNameUnit(ctx, &vocabulary)
		if err != nil {
			continue
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
