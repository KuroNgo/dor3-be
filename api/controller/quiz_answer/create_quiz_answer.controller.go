package quiz_answer_controller

import (
	quiz_answer_domain "clean-architecture/domain/quiz_answer"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (q *QuizAnswerController) CreateOneQuizAnswer(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	user, err := q.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var answerInput quiz_answer_domain.Input
	if err = ctx.ShouldBindJSON(&answerInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	answer := quiz_answer_domain.QuizAnswer{
		ID:          primitive.NewObjectID(),
		UserID:      user.ID,
		QuestionID:  answerInput.QuestionID,
		Answer:      answerInput.Answer,
		SubmittedAt: time.Now(),
	}

	err = q.QuizAnswerUseCase.CreateOne(ctx, &answer)
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
