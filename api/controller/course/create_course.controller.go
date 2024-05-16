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
	ctx.SSEvent("running", "Start create lesson management with file process")
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
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !file_internal.IsExcel(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not an Excel file",
		})
		return
	}

	ctx.SSEvent("running", "processing file")
	err = ctx.SaveUploadedFile(file, file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		err = os.Remove(file.Filename)
		if err != nil {
			fmt.Printf("Failed to delete temporary file: %v\n", err)
		}
	}()

	//resC, resL, resU, resV, err := excel.ReadFileForLessonManagementSystem(file.Filename)
	//if err != nil {
	//	_ = os.Remove(file.Filename)
	//	ctx.JSON(500, gin.H{"error": err.Error()})
	//	return
	//}

	resCourse, err := excel.ReadFileForCourse(file.Filename)
	if err != nil {
		_ = os.Remove(file.Filename)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.SSEvent("running", "Start create course")
	course := course_domain.Course{
		Id:          primitive.NewObjectID(),
		Name:        resCourse.Name,
		Description: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		WhoUpdated:  admin.FullName,
	}

	err = c.CourseUseCase.CreateOne(ctx, &course)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "create success course",
	})

	resLesson, err := excel.ReadFileForLesson(file.Filename)
	ctx.SSEvent("running", "Start create lesson")
	for _, data := range resLesson {
		courseID, errL := c.LessonUseCase.FindCourseIDByCourseName(ctx, data.CourseID)
		if errL != nil {
			_ = os.Remove(file.Filename)
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		l := lesson_domain.Lesson{
			ID:          primitive.NewObjectID(),
			CourseID:    courseID,
			Name:        data.Name,
			Level:       data.Level,
			IsCompleted: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			WhoUpdates:  admin.FullName,
		}

		// Tạo bài học trong cơ sở dữ liệu
		err = c.LessonUseCase.CreateOneByNameCourse(ctx, &l)
		if err != nil {
			_ = os.Remove(file.Filename)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "create success lesson",
	})

	ctx.SSEvent("running", "Start create unit")
	resUnit, err := excel.ReadFileForUnit(file.Filename)
	for _, unit := range resUnit {
		lessonID, errN := c.UnitUseCase.FindLessonIDByLessonName(ctx, unit.LessonID)
		if errN != nil {
			_ = os.Remove(file.Filename)
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		elUnit := unit_domain.Unit{
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

		err = c.UnitUseCase.CreateOneByNameLesson(ctx, &elUnit)
		if err != nil {
			continue
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "create success unit",
	})

	var vocabularies []vocabulary_domain.Vocabulary

	ctx.SSEvent("running", "Start create vocabulary")
	resVocabulary, err := excel.ReadFileForVocabulary(file.Filename)
	for _, vocabulary := range resVocabulary {
		unitID, errV := c.VocabularyUseCase.FindUnitIDByUnitLevel(ctx, vocabulary.UnitLevel, vocabulary.FieldOfIT)
		if errV != nil {
			_ = os.Remove(file.Filename)
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		vocabularyTrimSpace := strings.ReplaceAll(vocabulary.Word, " ", "")

		v := vocabulary_domain.Vocabulary{
			Id:            primitive.NewObjectID(),
			UnitID:        unitID,
			Word:          vocabulary.Word,
			WordForConfig: vocabularyTrimSpace,
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
		vocabularies = append(vocabularies, v)
	}

	for _, vocabulary := range vocabularies {
		err = c.VocabularyUseCase.CreateOneByNameUnit(ctx, &vocabulary)
		if err != nil {
			continue
		}
	}

	ctx.SSEvent("running", "Start create audio")
	data, err := c.VocabularyUseCase.GetLatestVocabulary(ctx)
	if err != nil {
		_ = os.Remove(file.Filename)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, audio := range data {
		_ = google.CreateTextToSpeech(audio)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success create audio for vocabulary latest",
	})

	ctx.SSEvent("running", "Start upload audio to cloud")
	// Lấy danh sách tệp trong thư mục audio
	dir := "audio"
	files, err := google.ListFilesInDirectory(dir)
	if err != nil {
		_ = os.Remove(file.Filename)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, audioFileName := range files {
		audio := strings.TrimSuffix(audioFileName, ".mp3")

		// Mở từng tệp
		f, errF := os.Open(filepath.Join(dir, audioFileName))
		if errF != nil {
			_ = os.Remove(file.Filename)
			return
		}

		// Kiểm tra xem file có phải là MP3 không
		if !file_internal.IsMP3(audioFileName) {
			_ = os.Remove(file.Filename)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("%s is not an MP3 file", audioFileName),
			})
			return
		}

		// Upload file lên Cloudinary
		dataRes, errD := cloudinary.UploadAudioToCloudinary(f, audioFileName, c.Database.CloudinaryUploadFolderAudioVocabulary)
		if errD != nil {
			_ = os.Remove(file.Filename)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Error uploading file %s to Cloudinary: %s", audioFileName, err),
			})
			return
		}

		// Tìm ID của từ vựng dựa trên tên file
		wordID, err := c.VocabularyUseCase.FindVocabularyIDByVocabularyName(ctx, audio)
		if err != nil {
			_ = os.Remove(file.Filename)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Error finding vocabulary ID for file %s: %s", audioFileName, err.Error()),
			})
			return
		}

		vocabularyRes := &vocabulary_domain.Vocabulary{
			Id:      wordID,
			LinkURL: dataRes.AudioURL,
		}

		errV := c.VocabularyUseCase.UpdateOneAudio(ctx, vocabularyRes)
		if errV != nil {
			_ = os.Remove(file.Filename)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Error updating vocabulary for file %s: %s", audioFileName, err.Error()),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success for upload audio to cloudinary ",
	})

	ctx.SSEvent("running", "Start delete audio")
	err = google.DeleteAllFilesInDirectory("audio")
	if err != nil {
		_ = os.Remove(file.Filename)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Trả về mảng dữ liệu dưới dạng JSON
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success for delete folder audio",
	})

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success create vocabulary",
	})

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success create vocabulary with file",
	})
}
