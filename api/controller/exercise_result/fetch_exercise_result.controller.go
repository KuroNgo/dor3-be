package exercise_result_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExerciseResultController) FetchManyExerciseResultInUser(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	user, err := e.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}
	page := ctx.DefaultQuery("page", "1")

	exam, err := e.ExerciseResultUseCase.FetchManyInUser(ctx, page)
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

func (e *ExerciseResultController) FetchResultByExerciseIDInUser(ctx *gin.Context) {
	exerciseID := ctx.Param("exercise_id")

	exercise, err := e.ExerciseResultUseCase.FetchManyByExerciseIDInUser(ctx, exerciseID)
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

func (e *ExerciseResultController) GetResultsExerciseIDInUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	exerciseId := ctx.Param("exercise_id")

	exercise, err := e.ExerciseResultUseCase.GetResultsExerciseIDInUser(ctx, fmt.Sprintf("%s", currentUser), exerciseId)
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
