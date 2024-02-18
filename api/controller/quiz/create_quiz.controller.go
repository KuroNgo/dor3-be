package quiz_controller

import (
	"clean-architecture/bootstrap"
	quiz_domain "clean-architecture/domain/quiz"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QuizCreateController struct {
	QuizUseCase quiz_domain.IQuizUseCase
	Database    *bootstrap.Database
}

func (quiz *QuizCreateController) CreateQuiz(ctx *gin.Context) {
	var quizInput quiz_domain.Input

	if err := ctx.ShouldBindJSON(&quizInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err := quiz.QuizUseCase.Create(ctx, &quizInput)
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
