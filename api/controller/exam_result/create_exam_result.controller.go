package exam_result_controller

import (
	exam_result_domain "clean-architecture/domain/exam_result"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (e *ExamResultController) CreateOneExam(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	userID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%s", currentUser))

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
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ExamID:    idExam,
		Score:     inputResult.Score,
		StartedAt: inputResult.StartedAt,
		Status:    1,
	}

	err := e.ExamResultUseCase.CreateOne(ctx, &result)
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
