package exercise_options_controller

import (
	exercise_options_domain "clean-architecture/domain/exercise_options"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (e *ExerciseOptionsController) UpdateOneExerciseOptions(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := e.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var answerInput exercise_options_domain.Input
	if err := ctx.ShouldBindJSON(&answerInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	answer := exercise_options_domain.ExerciseOptions{
		ID:         answerInput.ID,
		QuestionID: answerInput.QuestionID,
		Content:    answerInput.Content,
		BlankIndex: 0,
		UpdateAt:   time.Now(),
		WhoUpdate:  admin.FullName,
	}

	_, err = e.ExerciseOptionsUseCase.UpdateOne(ctx, &answer)
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