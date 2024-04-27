package quiz_result_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *QuizResultController) FetchManyExerciseResult(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	exam, err := e.QuizResultUseCase.FetchMany(ctx, page)
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

func (e *QuizResultController) FetchResultByQuizID(ctx *gin.Context) {
	exerciseID := ctx.Param("quiz_id")

	exercise, err := e.QuizResultUseCase.FetchManyByQuizID(ctx, exerciseID)
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

func (e *QuizResultController) GetResultsByUserIDAndQuizID(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	examId := ctx.Param("quiz_id")

	exam, err := e.QuizResultUseCase.GetResultsByUserIDAndQuizID(ctx, fmt.Sprintf("%s", currentUser), examId)
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
