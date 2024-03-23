package image_controller

import (
	image_domain "clean-architecture/domain/image"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (im *ImageController) CreateOneImageStatic(ctx *gin.Context) {
	// Parse form
	err := ctx.Request.ParseMultipartForm(4 << 20) // 8MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Just filename contains "file" in OS (operating system)
	file2, _, err := ctx.Request.FormFile("file")
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

	filename, ok := ctx.Get("filePath")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "filename not found",
		})
	}

	result, err := cloudinary.UploadToCloudinary(file2, filename.(string), im.Database.CloudinaryUploadFolderStatic)
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
		Category:  "static",
		AssetId:   result.AssetID,
	}

	// save data in database
	err = im.ImageUseCase.CreateOne(ctx, metadata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (im *ImageController) CreateOneImageLesson(ctx *gin.Context) {
	// Parse form
	err := ctx.Request.ParseMultipartForm(4 << 20) // 4 MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   "Error parsing form",
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Just filename contains "file" in OS (operating system)
	file2, _, err := ctx.Request.FormFile("file")
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

	filename, ok := ctx.Get("filePath")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "filename not found",
		})
	}

	result, err := cloudinary.UploadToCloudinary(file2, filename.(string), im.Database.CloudinaryUploadFolderLesson)
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
	err = im.ImageUseCase.CreateOne(ctx, metadata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (im *ImageController) CreateOneImageUser(ctx *gin.Context) {
	// Parse form
	err := ctx.Request.ParseMultipartForm(4 << 20) // 4 MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Just filename contains "file" in OS (operating system)
	file2, _, err := ctx.Request.FormFile("file")
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

	filename, ok := ctx.Get("filePath")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "filename not found",
		})
	}

	result, err := cloudinary.UploadToCloudinary(file2, filename.(string), im.Database.CloudinaryUploadFolderUser)
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
		Category:  "user",
		AssetId:   result.AssetID,
	}

	// save data in database
	err = im.ImageUseCase.CreateOne(ctx, metadata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (im *ImageController) CreateOneImageExam(ctx *gin.Context) {
	// Parse form
	err := ctx.Request.ParseMultipartForm(4 << 20) // 4 MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Just filename contains "file" in OS (operating system)
	file2, _, err := ctx.Request.FormFile("file")
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

	filename, ok := ctx.Get("filePath")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "filename not found",
		})
	}

	result, err := cloudinary.UploadToCloudinary(file2, filename.(string), im.Database.CloudinaryUploadFolderExam)
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
		Category:  "exam",
		AssetId:   result.AssetID,
	}

	// save data in database
	err = im.ImageUseCase.CreateOne(ctx, metadata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (im *ImageController) CreateOneImageQuiz(ctx *gin.Context) {
	// Parse form
	err := ctx.Request.ParseMultipartForm(4 << 20) // 4 MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Just filename contains "file" in OS (operating system)
	file2, _, err := ctx.Request.FormFile("file")
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

	filename, ok := ctx.Get("filePath")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "filename not found",
		})
	}

	result, err := cloudinary.UploadToCloudinary(file2, filename.(string), im.Database.CloudinaryUploadFolderQuiz)
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
		Category:  "quiz",
		AssetId:   result.AssetID,
	}

	// save data in database
	err = im.ImageUseCase.CreateOne(ctx, metadata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
