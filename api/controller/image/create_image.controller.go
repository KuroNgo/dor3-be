package image_controller

import (
	image_domain "clean-architecture/domain/image"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (im *ImageController) CreateImageInCloudinaryAndSaveMetadataInDatabase(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Just filename contains "file" in OS (operating system)
	file2, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !file_internal.IsImage(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename, ok := ctx.Get("filePath")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "filename not found"})
	}

	imageUrl, err := cloudinary.UploadToCloudinary(file2, filename.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// the metadata will be saved in database
	metadata := &image_domain.Image{
		Id:        primitive.NewObjectID(),
		ImageName: file.Filename,
		ImageUrl:  imageUrl,
		Size:      file.Size / 1024,
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
