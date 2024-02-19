package quiz_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (q *QuizController) FetchQuiz(ctx *gin.Context) {
	quiz, err := q.QuizUseCase.Fetch(ctx)
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
