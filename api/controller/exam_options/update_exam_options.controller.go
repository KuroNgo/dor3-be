package exam_options_controller

import (
	exam_options_domain "clean-architecture/domain/exam_options"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (e *ExamOptionsController) UpdateOneExamOptions(ctx *gin.Context) {
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

	var examInput exam_options_domain.Input
	if err := ctx.ShouldBindJSON(&examInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	exam := exam_options_domain.ExamOptions{
		ID:         examInput.ID,
		QuestionID: examInput.QuestionID,
		Content:    examInput.Content,
		UpdateAt:   time.Now(),
		WhoUpdate:  admin.FullName,
	}

	_, err = e.ExamOptionsUseCase.UpdateOne(ctx, &exam)
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
