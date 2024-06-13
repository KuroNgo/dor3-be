package exercise_controller

import (
	exercise_domain "clean-architecture/domain/exercise"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (e *ExerciseController) UpdateOneExercise(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	admin, err := e.AdminUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil || admin == nil {
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
		LessonID: exerciseInput.LessonID,
		UnitID:   exerciseInput.UnitID,

		Title:       exerciseInput.Title,
		Description: exerciseInput.Description,
		Duration:    exerciseInput.Duration,

		UpdatedAt:  time.Now(),
		WhoUpdates: admin.FullName,
	}

	data, err := e.ExerciseUseCase.UpdateOneInAdmin(ctx, &updateExercise)
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
