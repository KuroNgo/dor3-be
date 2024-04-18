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

	var quiz quiz_domain.Input
	if err := ctx.ShouldBindJSON(&quiz); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	updateQuiz := quiz_domain.Quiz{
		LessonID:     quiz.LessonID,
		UnitID:       quiz.UnitID,
		VocabularyID: quiz.VocabularyID,
		Title:        quiz.Title,
		Description:  quiz.Description,
		Duration:     quiz.Duration,
		WhoUpdates:   user.FullName,
		UpdatedAt:    time.Now(),
	}

	_, err = q.QuizUseCase.UpdateOne(ctx, &updateQuiz)
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

// UpdateOneQuiz done
func (q *QuizController) UpdateComplete(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")

	user, err := q.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		LessonID:     quiz.LessonID,
		UnitID:       quiz.UnitID,
		VocabularyID: quiz.VocabularyID,
		Title:        quiz.Title,
		Description:  quiz.Description,
		Duration:     quiz.Duration,
		WhoUpdates:   user.FullName,
		UpdatedAt:    time.Now(),
	}

	err = q.QuizUseCase.UpdateCompleted(ctx, &updateQuiz)
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
