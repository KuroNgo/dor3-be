package mark_vocabulary_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *MarkVocabularyController) DeleteOneMarkVocabulary(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	user, err := m.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	markVocabularyID := ctx.Query("_id")
	vocabularyID := ctx.Query("vocabulary_id")

	err = m.MarkVocabularyUseCase.DeleteOne(ctx, markVocabularyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err = m.VocabularyUseCase.UpdateIsFavouriteInUser(ctx, vocabularyID, 0)
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
