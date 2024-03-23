package mean_controller

import (
	mean_domain "clean-architecture/domain/mean"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (m *MeanController) CreateOneMean(ctx *gin.Context) {
	var meanInput mean_domain.Input
	if err := ctx.ShouldBindJSON(&meanInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := internal.IsValidMean(meanInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	meanRes := &mean_domain.Mean{
		ID:           primitive.NewObjectID(),
		VocabularyID: meanInput.VocabularyID,
		Description:  meanInput.Description,
		Example:      meanInput.Example,
		Synonym:      meanInput.Synonym,
		Antonym:      meanInput.Antonym,
	}

	err := m.MeanUseCase.CreateOne(ctx, meanRes)
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
