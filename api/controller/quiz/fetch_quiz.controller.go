package quiz_controller

import (
	"clean-architecture/bootstrap"
	quiz_domain "clean-architecture/domain/quiz"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QuizFetchController struct {
	QuizUseCase quiz_domain.IQuizUseCase
	Database    *bootstrap.Database
}

func (q *QuizFetchController) FetchQuiz(ctx *gin.Context) {
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
