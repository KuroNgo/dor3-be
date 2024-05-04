package image_controller

import (
	image_domain "clean-architecture/domain/image"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (im *ImageController) FetchImageByName(ctx *gin.Context) {
	var imageInput image_domain.Input
	if err := ctx.ShouldBindJSON(&imageInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	imageURL, err := im.ImageUseCase.GetURLByName(ctx, imageInput.ImageName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   imageURL,
	})
}

func (im *ImageController) FetchImageByCategory(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	category := ctx.Query("category")

	imageURL, err := im.ImageUseCase.FetchByCategory(ctx, category, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   imageURL,
	})
}

func (im *ImageController) FetchImage(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	imageURL, err := im.ImageUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   imageURL,
	})
}
