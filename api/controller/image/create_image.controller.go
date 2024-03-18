package image_controller

import (
	file_internal "clean-architecture/internal/file"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (im *ImageController) CreateImageInFireBaseAndSaveMetadataInDatabase(ctx *gin.Context) {
	// Just filename contains "file" in OS (operating system)
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !file_internal.IsImage(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Lưu file vào thư mục tạm thời
	err = ctx.SaveUploadedFile(file, "./"+file.Filename)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// the metadata will be saved in database
	//metadata := &image_domain.Image{
	//	Id:        primitive.NewObjectID(),
	//	ImageName: file.Filename,
	//	ImageUri:  ,
	//	Size:      file.Size,
	//}
}
