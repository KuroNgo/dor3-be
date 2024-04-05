package vocabulary_controller

import (
	file_internal "clean-architecture/internal/file"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Word struct {
	Vocabulary string `json:"word"`
}

func (v *VocabularyController) GenerateVoice(ctx *gin.Context) {
	var wordInput Word
	if err := ctx.ShouldBindJSON(&wordInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": "error",
		})
		return
	}

	_ = file_internal.CreateTextToSpeech(wordInput.Vocabulary)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "create success audio file",
	})
}
