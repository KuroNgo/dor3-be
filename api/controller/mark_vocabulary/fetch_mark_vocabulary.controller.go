package mark_vocabulary_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *MarkVocabularyController) FetchManyByMarkListIdAndUserId(ctx *gin.Context) {
	currentUser := ctx.MustGet("access_token")
	user, err := m.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))

	idMarkList := ctx.Query("course_id")

	markVocabulary, err := m.MarkVocabularyUseCase.FetchManyByMarkListIDAndUserId(ctx, idMarkList, fmt.Sprint(user.ID))
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
