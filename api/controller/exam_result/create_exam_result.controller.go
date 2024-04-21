package exam_result_controller

import (
	exam_result_domain "clean-architecture/domain/exam_result"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (e *ExamResultController) CreateOneExamResult(ctx *gin.Context) {
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

	examID := ctx.Query("exam_id")
	idExam, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%s", examID))

	var inputResult exam_result_domain.Input
	if err := ctx.ShouldBindJSON(&inputResult); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	result := exam_result_domain.ExamResult{
		ID:         primitive.NewObjectID(),
		UserID:     user.ID,
		ExamID:     idExam,
		Score:      inputResult.Score,
		StartedAt:  inputResult.StartedAt,
		IsComplete: 1,
	}

	err = e.ExamResultUseCase.CreateOne(ctx, &result)
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
