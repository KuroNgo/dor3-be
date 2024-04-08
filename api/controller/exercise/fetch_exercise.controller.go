package exercise_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExerciseController) FetchMany(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	
	lesson, err := e.ExerciseUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"lesson": lesson,
		},
	})
}
