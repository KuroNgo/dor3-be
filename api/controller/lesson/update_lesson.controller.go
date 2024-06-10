package lesson_controller

import (
	image_domain "clean-architecture/domain/image"
	lesson_domain "clean-architecture/domain/lesson"
	"clean-architecture/internal"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

// UpdateOneLesson handles the HTTP request to update a lesson.
// It verifies the current user's authorization and updates an existing lesson with the provided data.
func (l *LessonController) UpdateOneLesson(ctx *gin.Context) {
	// Retrieve the current user from the context
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		// Respond with unauthorized if no current user found
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  internal.FailStt,
			"message": internal.ErrUserNotLogin,
		})
		return
	}

	// Check if the current user is an admin
	admin, err := l.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		// Respond with unauthorized if the user is not an admin
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  internal.UnauthorizedStt,
			"message": internal.ErrUnauthorized,
		})
		return
	}

	// Bind JSON input to lessonInput struct
	var lessonInput lesson_domain.Input
	if err := ctx.ShouldBindJSON(&lessonInput); err != nil {
		// Respond with bad request if JSON binding fails
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Create a lesson struct with updated data
	updateLesson := lesson_domain.Lesson{
		ID:         lessonInput.ID,       // ID of the lesson to update
		CourseID:   lessonInput.CourseID, // Course ID the lesson belongs to
		Name:       lessonInput.Name,     // Updated lesson name
		Content:    lessonInput.Content,  // Updated lesson content
		UpdatedAt:  time.Now(),           // Set the updated timestamp
		WhoUpdates: admin.FullName,       // Set who performed the update
	}

	// Call the use case to update the lesson in the database
	data, err := l.LessonUseCase.UpdateOneInAdmin(ctx, &updateLesson)
	if err != nil {
		// Respond with bad request if update fails
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Respond with success and return the updated lesson data
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"result": data,
	})
}

// UpdateImageLesson handles the HTTP request to update a lesson with image
// It verifies the current user's authorization and updates an existing lesson with the provided data
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

	// the metadata will be saved in database
	metadata := &image_domain.Image{
		Id:        primitive.NewObjectID(),
		ImageName: file.Filename,
		Size:      file.Size / 1024,
		Category:  "lesson",
	}

	// save data in database
	err = l.ImageUseCase.CreateOne(ctx, metadata)
	if err != nil || err.Error() != internal.ErrExistDataInDatabase {
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
	dataUpdate := &image_domain.Image{
		Id:       primitive.NewObjectID(),
		ImageUrl: result.ImageURL,
		AssetId:  result.AssetID,
	}

	// save data in database
	err = l.ImageUseCase.UpdateOne(ctx, dataUpdate)
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

	data, err := l.LessonUseCase.UpdateOneInAdmin(ctx, &updateLesson)
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
