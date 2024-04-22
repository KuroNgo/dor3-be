package quiz_result_controller

import (
	quiz_result_domain "clean-architecture/domain/quiz_result"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (q *QuizResultController) CreateOneQuizResult(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	user, err := q.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	quizID := ctx.Query("quiz_id")
	idQuiz, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%s", quizID))

	var inputResult quiz_result_domain.Auto
	if err := ctx.ShouldBindJSON(&inputResult); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	result := quiz_result_domain.QuizResult{
		ID:         primitive.NewObjectID(),
		UserID:     user.ID,
		QuizID:     idQuiz,
		Score:      inputResult.Score,
		StartedAt:  inputResult.StartedAt,
		IsComplete: 1,
	}

	err = q.QuizResultUseCase.CreateOne(ctx, &result)
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
