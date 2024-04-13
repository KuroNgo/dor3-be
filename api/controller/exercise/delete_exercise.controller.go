package exercise_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExerciseController) DeleteOneExercise(ctx *gin.Context) {
	exerciseID := ctx.Query("_id")

	err := e.ExerciseUseCase.DeleteOne(ctx, exerciseID)
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
