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
	"strconv"
	"sync"
	"time"
)

func (l *LessonController) CreateOneLesson(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := l.AdminUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	courseId := ctx.Request.FormValue("course_id")
	idCourse, err := primitive.ObjectIDFromHex(courseId)

	name := ctx.Request.FormValue("name")
	content := ctx.Request.FormValue("content")
	le := ctx.Request.FormValue("level")
	level, _ := strconv.Atoi(le)

	file, err := ctx.FormFile("file")
	if err != nil {
		lessonRes := &lesson_domain.Lesson{
			ID:         primitive.NewObjectID(),
			CourseID:   idCourse,
			Name:       name,
			Content:    content,
			Level:      level,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			WhoUpdates: admin.FullName,
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
	result, err := cloudinary.UploadImageToCloudinary(f, file.Filename, l.Database.CloudinaryUploadFolderUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Tạo bài học với thông tin hình ảnh từ Cloudinary
	lessonRes := &lesson_domain.Lesson{
		ID:         primitive.NewObjectID(),
		CourseID:   idCourse,
		ImageURL:   result.ImageURL,
		Name:       name,
		Content:    content,
		Level:      level,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		WhoUpdates: admin.FullName,
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
}

func (l *LessonController) CreateOneLessonNotImage(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := l.AdminUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var inputLesson lesson_domain.Input
	if err := ctx.ShouldBindJSON(&inputLesson); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data",
		})
		return
	}

	lessonRes := &lesson_domain.Lesson{
		ID:         primitive.NewObjectID(),
		CourseID:   inputLesson.CourseID,
		Name:       inputLesson.Name,
		Content:    inputLesson.Content,
		Level:      inputLesson.Level,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		WhoUpdates: admin.FullName,
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
}

func (l *LessonController) CreateOneLessonHaveImage(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := l.AdminUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	courseID := ctx.Request.FormValue("course_id")
	idCourse, err := primitive.ObjectIDFromHex(courseID)

	name := ctx.Request.FormValue("name")
	content := ctx.Request.FormValue("content")
	le := ctx.Request.FormValue("level")
	level, _ := strconv.Atoi(le)

	file, err := ctx.FormFile("file")
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
	result, err := cloudinary.UploadImageToCloudinary(f, file.Filename, l.Database.CloudinaryUploadFolderUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Tạo bài học với thông tin hình ảnh từ Cloudinary
	lessonRes := &lesson_domain.Lesson{
		ID:         primitive.NewObjectID(),
		CourseID:   idCourse,
		ImageURL:   result.ImageURL,
		Name:       name,
		Content:    content,
		Level:      level,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		WhoUpdates: admin.FullName,
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
}

func (l *LessonController) CreateLessonWithFile(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := l.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
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

	err = ctx.SaveUploadedFile(file, file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save temporary file"})
		return
	}
	defer func() {
		err := os.Remove(file.Filename)
		if err != nil {
			fmt.Printf("Error: Failed to delete temporary file: %v\n", err)
		}
	}()

	result, err := excel.ReadFileForLesson(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read data from Excel file: " + err.Error()})
		return
	}

	var wg sync.WaitGroup
	var mux sync.Mutex
	errCh := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(errCh)

		for _, lesson := range result {
			// Tìm ID của khóa học từ tên khóa học
			mux.Lock()
			courseID, err := l.LessonUseCase.FindCourseIDByCourseName(ctx, lesson.CourseID)
			if err != nil {
				errCh <- err
				return
			}

			// Tạo bài học
			le := lesson_domain.Lesson{
				ID:          primitive.NewObjectID(),
				CourseID:    courseID,
				Name:        lesson.Name,
				Level:       lesson.Level,
				IsCompleted: 0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				WhoUpdates:  admin.FullName,
			}

			// Tạo bài học trong cơ sở dữ liệu
			err = l.LessonUseCase.CreateOneByNameCourse(ctx, &le)
			if err != nil {
				errCh <- err
				continue
			}
			mux.Unlock()
		}
	}()

	wg.Done()

	select {
	case err := <-errCh:
		ctx.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	default:
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}

}
