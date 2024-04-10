package mark_vocabulary_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *MarkVocabularyController) DeleteOneMarkVocabulary(ctx *gin.Context) {
	markVocabularyID := ctx.Query("_id")
	err := m.MarkVocabularyUseCase.DeleteOne(ctx, markVocabularyID)
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
