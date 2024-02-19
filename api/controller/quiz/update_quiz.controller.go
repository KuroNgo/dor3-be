package quiz_controller

import (
	quiz_domain "clean-architecture/domain/quiz"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (q *QuizController) UpdateQuiz(ctx *gin.Context) {
	var quiz quiz_domain.Input
	if err := ctx.ShouldBindJSON(&quiz); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	updateQuiz := quiz_domain.Quiz{
		Question:      quiz.Question,
		Options:       quiz.Options,
		CorrectAnswer: quiz.CorrectAnswer,
		QuestionType:  quiz.QuestionType,
	}

	err := q.QuizUseCase.Update(ctx, updateQuiz.ID.Hex(), updateQuiz)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
