package quiz_options_controller

import (
	quiz_options_domain "clean-architecture/domain/quiz_options"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (e *QuizOptionsController) UpdateOneQuizOptions(ctx *gin.Context) {
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

	var answerInput quiz_options_domain.Input
	if err := ctx.ShouldBindJSON(&answerInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	answer := quiz_options_domain.QuizOptions{
		ID:         answerInput.ID,
		QuestionID: answerInput.QuestionID,
		Content:    answerInput.Content,
		UpdateAt:   time.Now(),
		WhoUpdate:  admin.FullName,
	}

	_, err = e.QuizOptionsUseCase.UpdateOne(ctx, &answer)
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
