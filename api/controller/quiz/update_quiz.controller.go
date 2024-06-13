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
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	admin, err := q.AdminUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var quiz quiz_domain.Input
	if err := ctx.ShouldBindJSON(&quiz); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	updateQuiz := quiz_domain.Quiz{
		LessonID:    quiz.LessonID,
		UnitID:      quiz.UnitID,
		Title:       quiz.Title,
		Description: quiz.Description,
		Duration:    quiz.Duration,
		WhoUpdates:  admin.FullName,
		UpdatedAt:   time.Now(),
	}

	_, err = q.QuizUseCase.UpdateOneInAdmin(ctx, &updateQuiz)
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
