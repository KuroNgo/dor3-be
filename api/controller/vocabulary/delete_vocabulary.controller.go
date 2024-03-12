package vocabulary_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (v *VocabularyController) DeleteOneVocabulary(ctx *gin.Context) {
	vocabularyID := ctx.Query("_id")

	err := v.VocabularyUseCase.DeleteOne(ctx, vocabularyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Trả về mảng dữ liệu dưới dạng JSON
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "the vocabulary is deleted!",
	})
}
