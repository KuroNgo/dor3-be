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
	admin, err := e.AdminUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

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
		Id:           primitive.NewObjectID(),
		LessonID:     exerciseInput.LessonID,
		UnitID:       exerciseInput.UnitID,
		VocabularyID: exerciseInput.VocabularyID,

		Title:       exerciseInput.Title,
		Description: exerciseInput.Description,
		Duration:    exerciseInput.Duration,

		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		WhoUpdates: admin.FullName,
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
