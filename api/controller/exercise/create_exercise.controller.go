package exercise_controller

import (
	exercise_domain "clean-architecture/domain/exercise"
	"clean-architecture/internal"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (e *ExerciseController) CreateOneExercise(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	user, err := e.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))

	var exerciseInput exercise_domain.Input
	if err := ctx.ShouldBindJSON(&exerciseInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if err := internal.IsValidExercise(exerciseInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	exerciseRes := &exercise_domain.Exercise{
		Id:         primitive.NewObjectID(),
		Vocabulary: exerciseInput.VocabularyID,
		Title:      exerciseInput.Title,
		Content:    exerciseInput.Content,
		Type:       exerciseInput.Type,
		//Options:    exerciseInput.Options,
		CorrectAns: exerciseInput.CorrectAns,
		BlankIndex: exerciseInput.BlankIndex,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		WhoUpdates: user.FullName,
	}

	err = e.ExerciseUseCase.CreateOne(ctx, exerciseRes)
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
