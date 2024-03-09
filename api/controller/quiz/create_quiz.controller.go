package quiz_controller

import (
	quiz_domain "clean-architecture/domain/quiz"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (q *QuizController) CreateOneQuiz(ctx *gin.Context) {
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

	quizRes := &quiz_domain.Quiz{
		ID:            primitive.NewObjectID(),
		Question:      quizInput.Question,
		Options:       quizInput.Options,
		CorrectAnswer: quizInput.CorrectAnswer,
		Explanation:   quizInput.Explanation,
		QuestionType:  quizInput.QuestionType,
	}

	err := q.QuizUseCase.CreateOne(ctx, quizRes)
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
