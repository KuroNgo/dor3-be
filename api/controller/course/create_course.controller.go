package course_controller

import (
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/internal"
	"clean-architecture/internal/cloud/cloudinary"
	"clean-architecture/internal/cloud/google"
	file_internal "clean-architecture/internal/file"
	"clean-architecture/internal/file/excel"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *CourseController) CreateOneCourseInAdmin(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := c.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

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
		WhoUpdated:  admin.FullName,
	}

	err = c.CourseUseCase.CreateOneInAdmin(ctx, course)
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

func (c *CourseController) CreateCourseWithFileInAdmin(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := c.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	err = ctx.Request.ParseMultipartForm(8 << 20) // 8MB max size
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
		_ = os.Remove(file.Filename)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	course := course_domain.Course{
		Id:          primitive.NewObjectID(),
		Name:        result.Name,
		Description: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		WhoUpdated:  admin.FullName,
	}

	err = c.CourseUseCase.CreateOneInAdmin(ctx, &course)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		_ = os.Remove(file.Filename)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

}

func (c *CourseController) CreateLessonManagementWithFileInAdmin(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := c.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	err = ctx.Request.ParseMultipartForm(8 << 20) // 8MB max size
	if err != nil {
		ctx.String(http.StatusBadRequest, "Error parsing form: "+err.Error())
		return
	}

	file, err := ctx.FormFile("files")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !file_internal.IsExcel(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Not an Excel file"})
		return
	}

	err = ctx.SaveUploadedFile(file, file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if err := os.Remove(file.Filename); err != nil {
			fmt.Printf("Failed to delete temporary file: %v\n", err)
		}
	}()

	resCourse, err := excel.ReadFileForCourse(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	course := course_domain.Course{
		Id:          primitive.NewObjectID(),
		Name:        resCourse.Name,
		Description: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		WhoUpdated:  admin.FullName,
	}

	_ = c.CourseUseCase.CreateOneInAdmin(ctx, &course)

	resLesson, err := excel.ReadFileForLesson(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, data := range resLesson {
		courseID, err := c.CourseUseCase.FindCourseIDByCourseNameInAdmin(ctx, data.CourseID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lesson := lesson_domain.Lesson{
			ID:         primitive.NewObjectID(),
			CourseID:   courseID,
			Name:       data.Name,
			Level:      data.Level,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			WhoUpdates: admin.FullName,
		}

		_ = c.LessonUseCase.CreateOneByNameCourseInAdmin(ctx, &lesson)
	}

	resUnit, err := excel.ReadFileForUnit(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, unit := range resUnit {
		lessonID, err := c.LessonUseCase.FindLessonIDByLessonNameInAdmin(ctx, unit.LessonID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		un := unit_domain.Unit{
			ID:        primitive.NewObjectID(),
			LessonID:  lessonID,
			Name:      unit.Name,
			Level:     unit.Level,
			ImageURL:  "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			WhoCreate: admin.FullName,
		}

		_ = c.UnitUseCase.CreateOneByNameLessonInAdmin(ctx, &un)
	}

	resVocabulary, err := excel.ReadFileForVocabulary(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var vocabularies []vocabulary_domain.Vocabulary
	for _, vocabulary := range resVocabulary {
		unitID, err := c.VocabularyUseCase.FindUnitIDByUnitLevelInAdmin(ctx, vocabulary.UnitLevel, vocabulary.FieldOfIT)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		vocab := vocabulary_domain.Vocabulary{
			Id:            primitive.NewObjectID(),
			UnitID:        unitID,
			Word:          vocabulary.Word,
			WordForConfig: strings.ReplaceAll(vocabulary.Word, " ", ""),
			PartOfSpeech:  vocabulary.PartOfSpeech,
			Pronunciation: vocabulary.Pronunciation,
			Mean:          vocabulary.Example,
			ExampleEng:    vocabulary.ExampleEng,
			ExampleVie:    vocabulary.ExampleVie,
			ExplainEng:    vocabulary.ExplainEng,
			ExplainVie:    vocabulary.ExplainVie,
			FieldOfIT:     vocabulary.FieldOfIT,
			LinkURL:       "",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			WhoUpdates:    admin.FullName,
		}
		vocabularies = append(vocabularies, vocab)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	vocabulariesCount := 0
	go func() {
		defer wg.Done()
		for _, vocabulary := range vocabularies {
			if err = c.VocabularyUseCase.CreateOneByNameUnitInAdmin(ctx, &vocabulary); err != nil {
				continue
			}
			vocabulariesCount++
		}
	}()
	wg.Wait()

	if vocabulariesCount > 0 {
		data, err := c.VocabularyUseCase.GetLatestVocabularyInAdmin(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, audio := range data {
				if err = google.CreateTextToSpeech(audio); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}()
		wg.Wait()

		dir := "audio"
		files, err := google.ListFilesInDirectory(dir)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, audioFileName := range files {
				audio := strings.TrimSuffix(audioFileName, ".mp3")

				f, err := os.Open(filepath.Join(dir, audioFileName))
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				defer f.Close()

				if !file_internal.IsMP3(audioFileName) {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is not an MP3 file", audioFileName)})
					return
				}

				dataRes, err := cloudinary.UploadAudioToCloudinary(f, audioFileName, c.Database.CloudinaryUploadFolderAudioVocabulary)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error uploading file %s to Cloudinary: %s", audioFileName, err)})
					return
				}

				wordID, err := c.VocabularyUseCase.FindVocabularyIDByVocabularyConfigInAdmin(ctx, audio)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error finding vocabulary ID for file %s: %s", audioFileName, err)})
					return
				}

				vocabularyRes := &vocabulary_domain.Vocabulary{
					Id:      wordID,
					LinkURL: dataRes.AudioURL,
				}

				if err := c.VocabularyUseCase.UpdateOneAudioInAdmin(ctx, vocabularyRes); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating vocabulary for file %s: %s", audioFileName, err)})
					return
				}
			}
		}()
		wg.Wait()

		if err = google.DeleteAllFilesInDirectory("audio"); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "success create vocabulary with file"})
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"status":  "error",
		"message": "The vocabulary in the files already exists.",
	})
}
