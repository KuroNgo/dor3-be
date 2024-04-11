package mark_vocabulary_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *MarkVocabularyController) DeleteOneMarkVocabulary(ctx *gin.Context) {
	markVocabularyID := ctx.Query("_id")
	vocabularyID := ctx.Query("vocabulary_id")

	err := m.MarkVocabularyUseCase.DeleteOne(ctx, markVocabularyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err = m.VocabularyUseCase.UpdateIsFavourite(ctx, vocabularyID, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error create mark vocabulary",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
