package quiz_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// FetchManyQuiz done
func (q *QuizController) FetchManyQuiz(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	page := ctx.DefaultQuery("page", "1")
	quiz, err := q.QuizUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"quiz": quiz,
		},
	})
}
