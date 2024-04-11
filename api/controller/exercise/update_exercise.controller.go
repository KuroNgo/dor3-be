package exercise_controller

import (
	exercise_domain "clean-architecture/domain/exercise"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (e *ExerciseController) UpdateOneExercise(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")

	user, err := e.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	exerciseID := ctx.Query("_id")

	var exerciseInput exercise_domain.Input
	if err := ctx.ShouldBindJSON(&exerciseInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	updateExercise := exercise_domain.Exercise{
		Vocabulary: exerciseInput.VocabularyID,
		Title:      exerciseInput.Title,
		Content:    exerciseInput.Content,
		Type:       exerciseInput.Type,
		//Options:    exerciseInput.Options,
		CorrectAns: exerciseInput.CorrectAns,
		BlankIndex: exerciseInput.BlankIndex,
		UpdatedAt:  time.Now(),
		WhoUpdates: user.FullName,
	}

	err = e.ExerciseUseCase.UpdateOne(ctx, exerciseID, updateExercise)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

}
