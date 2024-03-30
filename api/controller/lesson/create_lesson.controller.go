package lesson_controller

import (
	lesson_domain "clean-architecture/domain/lesson"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (l *LessonController) CreateOneLesson(ctx *gin.Context) {
	// Lấy user hiện tại đăng nhập
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

	// Nếu có file được tải lên, tiến hành xử lý file
	file, err = ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	// Kiểm tra xem file có phải là hình ảnh không
	if !file_internal.IsImage(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file format. Only images are allowed.",
		})
		return
	}

	// Mở file để đọc dữ liệu
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
