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

	// Parse form
	err = ctx.Request.ParseMultipartForm(8 << 20) // 8MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	// Get uploaded file
	file, err := ctx.FormFile("files")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kiểm tra định dạng file
	if !file_internal.IsExcel(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not an Excel file",
		})
		return
	}

	// Lưu file vào thư mục tạm
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

	// Đọc dữ liệu từ file Excel
	result, err := excel.ReadFileForLesson(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Tạo kênh để gửi và nhận lỗi từ goroutine
	errChan := make(chan error)
	defer close(errChan)

	// Tạo các bài học từ dữ liệu bằng cách sử dụng goroutine
	successCount := 0
	for _, lesson := range result {
		go func(lesson file_internal.Lesson) {
			// Tìm ID của khóa học từ tên khóa học
			courseID, err := l.LessonUseCase.FindCourseIDByCourseName(ctx, lesson.CourseID)
			if err != nil {
				errChan <- fmt.Errorf("failed to find course ID for course '%s': %v", lesson.CourseID, err)
				return
			}

			// Tạo bài học
			le := lesson_domain.Lesson{
				ID:          primitive.NewObjectID(),
				CourseID:    courseID,
				Name:        lesson.Name,
				Content:     lesson.Content,
				Level:       lesson.Level,
				IsCompleted: 0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				WhoUpdates:  user.FullName,
			}

			// Tạo bài học trong cơ sở dữ liệu
			err = l.LessonUseCase.CreateOneByNameCourse(ctx, &le)
			if err != nil {
				errChan <- fmt.Errorf("failed to create lesson '%s': %v", le.Name, err)
				return
			}

			// Gửi thông báo thành công
			errChan <- nil
		}(lesson)
	}

	// Đợi goroutine kết thúc và xử lý lỗi nếu có
	for range result {
		if err := <-errChan; err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			successCount++
		}
	}

	// Trả về kết quả
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
