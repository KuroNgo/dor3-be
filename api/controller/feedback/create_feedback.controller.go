package feedback_controller

import (
	feedback_domain "clean-architecture/domain/feedback"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (f *FeedbackController) CreateOneFeedback(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	user, err := f.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var feedbackInput feedback_domain.Input
	if err := ctx.ShouldBindJSON(&feedbackInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	feedback := feedback_domain.Feedback{
		ID:            primitive.NewObjectID(),
		UserID:        user.ID,
		Title:         feedbackInput.Title,
		Content:       feedbackInput.Content,
		Feeling:       feedbackInput.Feeling,
		SubmittedDate: time.Now(),
		IsLoveWeb:     feedbackInput.IsLoveWeb,
		IsSeen:        0,
	}

	err = f.FeedbackUseCase.CreateOneByUser(ctx, &feedback)
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
