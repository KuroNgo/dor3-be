package quiz_question_controller

import (
	quiz_question_domain "clean-architecture/domain/quiz_question"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (q *QuizQuestionsController) UpdateOneQuizOptions(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := q.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var questionInput quiz_question_domain.Input
	if err := ctx.ShouldBindJSON(&questionInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	questions := quiz_question_domain.QuizQuestion{
		ID:           primitive.NewObjectID(),
		QuizID:       questionInput.QuizID,
		VocabularyID: questionInput.VocabularyID,
		Content:      questionInput.Content,
		Type:         questionInput.Type,
		Level:        questionInput.Level,
		UpdateAt:     time.Now(),
		WhoUpdate:    admin.FullName,
	}

	_, err = q.QuizQuestionUseCase.UpdateOne(ctx, &questions)
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
