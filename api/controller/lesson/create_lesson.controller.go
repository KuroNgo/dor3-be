package lesson_controller

import (
	lesson_domain "clean-architecture/domain/lesson"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"clean-architecture/internal/file/excel"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"time"
)

func (l *LessonController) CreateOneLesson(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")

	user, err := l.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	var lessonInput lesson_domain.Input
	if err = ctx.ShouldBindJSON(&lessonInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		lessonRes := &lesson_domain.Lesson{
			ID:         primitive.NewObjectID(),
			CourseID:   lessonInput.CourseID,
			Name:       lessonInput.Name,
			Content:    lessonInput.Content,
			Level:      lessonInput.Level,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			WhoUpdates: user.FullName,
		}
		err = l.LessonUseCase.CreateOne(ctx, lessonRes)
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
		return
	}

	file, err = ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	if !file_internal.IsImage(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file format. Only images are allowed.",
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error opening uploaded file",
		})
		return
	}
	defer f.Close()

	// Tải file lên Cloudinary
	filename, _ := ctx.Get("filePath")
	result, err := cloudinary.UploadToCloudinary(f, filename.(string), l.Database.CloudinaryUploadFolderUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Tạo bài học với thông tin hình ảnh từ Cloudinary
	lessonRes := &lesson_domain.Lesson{
		ID:         primitive.NewObjectID(),
		CourseID:   lessonInput.CourseID,
		ImageURL:   result.ImageURL,
		Name:       lessonInput.Name,
		Content:    lessonInput.Content,
		Level:      lessonInput.Level,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		WhoUpdates: user.FullName,
	}

	err = l.LessonUseCase.CreateOne(ctx, lessonRes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Trả về kết quả thành công nếu không có lỗi xảy ra
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (l *LessonController) CreateLessonWithFile(ctx *gin.Context) {
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

	result, err := excel.ReadFileForLesson(file.Filename)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var lessons []lesson_domain.Lesson
	for _, lesson := range result {
		courseID, err := l.LessonUseCase.FindCourseIDByCourseName(ctx, lesson.CourseID)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		l := lesson_domain.Lesson{
			ID:          primitive.NewObjectID(),
			CourseID:    courseID,
			Name:        lesson.Name,
			Content:     lesson.Content,
			Level:       lesson.Level,
			IsCompleted: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			//WhoUpdates:
		}
		lessons = append(lessons, l)
	}

	successCount := 0
	for _, lesson := range lessons {
		err = l.LessonUseCase.CreateOneByNameCourse(ctx, &lesson)
		if err != nil {
			continue
		}
		successCount++
	}

	if successCount == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create any lesson",
			"message": "Any value have exist in database",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"success_count": successCount,
	})
}
