package exam_question_controller

import (
	exam_question_domain "clean-architecture/domain/exam_question"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (e *ExamQuestionsController) UpdateOneExamQuestion(ctx *gin.Context) {
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

	var questionInput exam_question_domain.Input
	if err := ctx.ShouldBindJSON(&questionInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	questions := exam_question_domain.ExamQuestion{
		ID:           questionInput.ID,
		ExamID:       questionInput.ExamID,
		VocabularyID: questionInput.VocabularyID,
		Content:      questionInput.Content,
		Type:         questionInput.Type,
		Level:        questionInput.Level,
		UpdateAt:     time.Now(),
		WhoUpdate:    admin.FullName,
	}

	data, err := e.ExamQuestionUseCase.UpdateOne(ctx, &questions)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
