package vocabulary_controller

import (
	"clean-architecture/internal/cloud/google"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Word struct {
	Vocabulary string `json:"word"`
}

func (v *VocabularyController) GenerateVoice(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := v.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var wordInput Word
	if err := ctx.ShouldBindJSON(&wordInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": "error",
		})
		return
	}

	err = google.CreateTextToSpeech(wordInput.Vocabulary)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "create success audio file",
	})
}
