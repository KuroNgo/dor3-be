package exercise_result_controller

import (
	exercise_result_domain "clean-architecture/domain/exercise_result"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (e *ExerciseResultController) CreateOneExercise(ctx *gin.Context) {
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

	exerciseID := ctx.Query("exercise_id")
	idExercise, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%s", exerciseID))

	var inputResult exercise_result_domain.Input
	if err := ctx.ShouldBindJSON(&inputResult); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	result := exercise_result_domain.ExerciseResult{
		ID:         primitive.NewObjectID(),
		UserID:     user.ID,
		ExerciseID: idExercise,
		Score:      inputResult.Score,
		StartedAt:  inputResult.StartedAt,
		IsComplete: 1,
	}

	err = e.ExerciseResultUseCase.CreateOneInUser(ctx, &result)
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
