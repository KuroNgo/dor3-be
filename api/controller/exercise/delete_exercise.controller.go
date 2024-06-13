package exercise_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExerciseController) DeleteOneExercise(ctx *gin.Context) {
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
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": admin.FullName + "You are not authorized to perform this action!",
		})
		return
	}
	exerciseID := ctx.Query("_id")

	err = e.ExerciseUseCase.DeleteOneInAdmin(ctx, exerciseID)
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
