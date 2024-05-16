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
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (c *CourseController) CreateOneCourse(ctx *gin.Context) {
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

	err = c.CourseUseCase.CreateOne(ctx, course)
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
	currentUser := ctx.MustGet("currentUser")
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

	err = c.CourseUseCase.CreateOne(ctx, &course)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error creating course: %v\n": err,
		})
		_ = os.Remove(file.Filename)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

}

func (c *CourseController) CreateLessonManagementWithFile(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Flush()

	sendSSE := func(event, message string) {
		ctx.SSEvent(event, message)
		ctx.Writer.Flush()
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			sendSSE("keep-alive", "ping")
		}
	}()

	sendSSE("running", "Start create lesson management with file process")

	currentUser := ctx.MustGet("currentUser")
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

	sendSSE("running", "Processing file")
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

	sendSSE("running", "Start create course")
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

	if err := c.CourseUseCase.CreateOne(ctx, &course); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sendSSE("running", "Start create lesson")
	resLesson, err := excel.ReadFileForLesson(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, data := range resLesson {
		courseID, err := c.LessonUseCase.FindCourseIDByCourseName(ctx, data.CourseID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lesson := lesson_domain.Lesson{
			ID:          primitive.NewObjectID(),
			CourseID:    courseID,
			Name:        data.Name,
			Level:       data.Level,
			IsCompleted: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			WhoUpdates:  admin.FullName,
		}

		if err := c.LessonUseCase.CreateOneByNameCourse(ctx, &lesson); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	sendSSE("running", "Start create unit")
	resUnit, err := excel.ReadFileForUnit(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, unit := range resUnit {
		lessonID, err := c.UnitUseCase.FindLessonIDByLessonName(ctx, unit.LessonID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		unit := unit_domain.Unit{
			ID:         primitive.NewObjectID(),
			LessonID:   lessonID,
			Name:       unit.Name,
			Level:      unit.Level,
			ImageURL:   "",
			IsComplete: 0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			WhoCreate:  admin.FullName,
		}

		if err := c.UnitUseCase.CreateOneByNameLesson(ctx, &unit); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	sendSSE("running", "Start create vocabulary")
	resVocabulary, err := excel.ReadFileForVocabulary(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var vocabularies []vocabulary_domain.Vocabulary
	for _, vocabulary := range resVocabulary {
		unitID, err := c.VocabularyUseCase.FindUnitIDByUnitLevel(ctx, vocabulary.UnitLevel, vocabulary.FieldOfIT)
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
			IsFavourite:   0,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			WhoUpdates:    admin.FullName,
		}
		vocabularies = append(vocabularies, vocab)
	}

	for _, vocabulary := range vocabularies {
		if err := c.VocabularyUseCase.CreateOneByNameUnit(ctx, &vocabulary); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	sendSSE("running", "Start create audio")
	data, err := c.VocabularyUseCase.GetLatestVocabulary(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, audio := range data {
		if err := google.CreateTextToSpeech(audio); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	sendSSE("running", "Start upload audio to cloud")
	dir := "audio"
	files, err := google.ListFilesInDirectory(dir)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

		wordID, err := c.VocabularyUseCase.FindVocabularyIDByVocabularyName(ctx, audio)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error finding vocabulary ID for file %s: %s", audioFileName, err)})
			return
		}

		vocabularyRes := &vocabulary_domain.Vocabulary{
			Id:      wordID,
			LinkURL: dataRes.AudioURL,
		}

		if err := c.VocabularyUseCase.UpdateOneAudio(ctx, vocabularyRes); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating vocabulary for file %s: %s", audioFileName, err)})
			return
		}
	}

	sendSSE("running", "Start delete audio")
	if err := google.DeleteAllFilesInDirectory("audio"); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success create vocabulary with file"})
}
