package image_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (im *ImageController) FetchImageByName(ctx *gin.Context) {
	imageName := ctx.Query("name")

	imageURL, err := im.ImageUseCase.GetURLByName(ctx, imageName)
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
