package mark_vocabulary_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *MarkVocabularyController) FetchManyByMarkListIdAndUserId(ctx *gin.Context) {
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

	idMarkList := ctx.Query("mark_list_id")
	markVocabulary, err := m.MarkVocabularyUseCase.FetchManyByMarkListIDAndUserId(ctx, idMarkList, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"mark_vocabulary": markVocabulary,
	})
}

func (m *MarkVocabularyController) FetchManyByMarkListIdInAdmin(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	admin, err := m.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	idMarkList := ctx.Query("mark_list_id")
	markVocabulary, err := m.MarkVocabularyUseCase.FetchManyByMarkList(ctx, idMarkList)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"mark_vocabulary": markVocabulary,
	})
}
