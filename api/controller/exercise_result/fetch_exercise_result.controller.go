package exercise_result_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExerciseResultController) FetchManyExerciseResult(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	exam, err := e.ExerciseResultUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   exam,
	})
}

func (e *ExerciseResultController) FetchResultByExerciseID(ctx *gin.Context) {
	exerciseID := ctx.Query("exercise_id")

	exercise, err := e.ExerciseResultUseCase.FetchManyByExerciseID(ctx, exerciseID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   exercise,
	})
}

func (e *ExerciseResultController) GetResultsByUserIDAndExerciseID(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	exerciseId := ctx.Query("exercise_id")

	exercise, err := e.ExerciseResultUseCase.GetResultsByUserIDAndExerciseID(ctx, fmt.Sprintf("%s", currentUser), exerciseId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   exercise,
	})
}
