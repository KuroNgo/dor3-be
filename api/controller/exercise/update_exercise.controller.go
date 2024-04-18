package exercise_controller

import (
	exercise_domain "clean-architecture/domain/exercise"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (e *ExerciseController) UpdateOneExercise(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")

	user, err := e.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var exerciseInput exercise_domain.Input
	if err := ctx.ShouldBindJSON(&exerciseInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	updateExercise := exercise_domain.Exercise{
		LessonID:     exerciseInput.LessonID,
		UnitID:       exerciseInput.UnitID,
		VocabularyID: exerciseInput.VocabularyID,

		Title:       exerciseInput.Title,
		Description: exerciseInput.Description,
		Duration:    exerciseInput.Duration,

		UpdatedAt:  time.Now(),
		WhoUpdates: user.FullName,
	}

	data, err := e.ExerciseUseCase.UpdateOne(ctx, &updateExercise)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})

}
