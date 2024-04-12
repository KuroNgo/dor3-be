package exercise_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExerciseController) FetchMany(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	exercise, err := e.ExerciseUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   exercise,
	})
}
