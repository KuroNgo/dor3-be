package quiz_controller

import (
	quiz_domain "clean-architecture/domain/quiz"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// UpdateOneQuiz done
func (q *QuizController) UpdateOneQuiz(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")

	user, err := q.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quizID := ctx.Query("_id")

	var quiz quiz_domain.Input
	if err := ctx.ShouldBindJSON(&quiz); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	updateQuiz := quiz_domain.Quiz{
		WhoUpdates: user.FullName,
		UpdatedAt:  time.Now(),
	}

	err = q.QuizUseCase.UpdateOne(ctx, quizID, updateQuiz)
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
