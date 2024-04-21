package exercise_quesiton_controller

import (
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (e *ExerciseQuestionsController) CreateOneExerciseQuestions(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := e.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var questionInput exercise_questions_domain.Input
	if err = ctx.ShouldBindJSON(&questionInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	question := exercise_questions_domain.ExerciseQuestion{
		ID:           primitive.NewObjectID(),
		ExerciseID:   questionInput.ExerciseID,
		VocabularyID: questionInput.VocabularyID,
		Content:      questionInput.Content,
		Type:         questionInput.Type,
		Level:        questionInput.Level,
		CreatedAt:    time.Now(),
		UpdateAt:     time.Now(),
		WhoUpdate:    admin.FullName,
	}

	err = e.ExerciseQuestionUseCase.CreateOne(ctx, &question)
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
