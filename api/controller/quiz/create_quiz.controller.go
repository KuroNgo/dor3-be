package quiz_controller

import (
	quiz_domain "clean-architecture/domain/quiz"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (q *QuizController) CreateQuiz(ctx *gin.Context) {
	var quizInput quiz_domain.Input

	if err := ctx.ShouldBindJSON(&quizInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := internal.IsValidQuiz(quizInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err := q.QuizUseCase.Create(ctx, &quizInput)
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
