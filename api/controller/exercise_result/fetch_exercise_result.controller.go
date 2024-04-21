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

	exercise, err := e.ExerciseResultUseCase.FetchManyByExamID(ctx, exerciseID)
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

func (e *ExerciseResultController) GetResultsByUserIDAndExamID(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	examId := ctx.Query("exam_id")

	exam, err := e.ExerciseResultUseCase.GetResultsByUserIDAndExamID(ctx, fmt.Sprintf("%s", currentUser), examId)
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
