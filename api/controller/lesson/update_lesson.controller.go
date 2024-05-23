package lesson_controller

import (
	image_domain "clean-architecture/domain/image"
	lesson_domain "clean-architecture/domain/lesson"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (l *LessonController) UpdateOneLesson(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := l.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var lessonInput lesson_domain.Input
	if err := ctx.ShouldBindJSON(&lessonInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	updateLesson := lesson_domain.Lesson{
		ID:         lessonInput.ID,
		CourseID:   lessonInput.CourseID,
		Name:       lessonInput.Name,
		Content:    lessonInput.Content,
		UpdatedAt:  time.Now(),
		WhoUpdates: admin.FullName,
	}

	data, err := l.LessonUseCase.UpdateOne(ctx, &updateLesson)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"result": data,
	})
}

func (l *LessonController) UpdateImageLesson(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := l.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	// Parse form
	err = ctx.Request.ParseMultipartForm(4 << 20) // 8MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("files")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !file_internal.IsImage(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := cloudinary.UploadImageToCloudinary(f, file.Filename, l.Database.CloudinaryUploadFolderStatic)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// the metadata will be saved in database
	metadata := &image_domain.Image{
		Id:        primitive.NewObjectID(),
		ImageName: file.Filename,
		ImageUrl:  result.ImageURL,
		Size:      file.Size / 1024,
		Category:  "lesson",
		AssetId:   result.AssetID,
	}

	// save data in database
	err = l.ImageUseCase.CreateOne(ctx, metadata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	lessonId := ctx.Request.FormValue("_id")
	idLesson, _ := primitive.ObjectIDFromHex(lessonId)

	updateLesson := lesson_domain.Lesson{
		ID:         idLesson,
		ImageURL:   result.ImageURL,
		AssetURL:   result.AssetID,
		UpdatedAt:  time.Now(),
		WhoUpdates: admin.FullName,
	}

	data, err := l.LessonUseCase.UpdateOne(ctx, &updateLesson)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"result": data,
	})
}
