package quiz_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// FetchManyQuiz done
func (q *QuizController) FetchManyQuiz(ctx *gin.Context) {
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
