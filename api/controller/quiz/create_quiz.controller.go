package quiz_controller

import (
	quiz_domain "clean-architecture/domain/quiz"
	"clean-architecture/internal"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (q *QuizController) CreateOneQuiz(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")

	user, err := q.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		WhoUpdates:    user.FullName,
	}

	err = q.QuizUseCase.CreateOne(ctx, quizRes)
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

func (q *QuizController) CreateFromFileQuiz(ctx *gin.Context) {

}
