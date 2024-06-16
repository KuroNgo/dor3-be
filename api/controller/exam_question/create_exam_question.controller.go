package exam_question_controller

import (
	exam_question_domain "clean-architecture/domain/exam_question"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (e *ExamQuestionsController) CreateOneExamQuestions(ctx *gin.Context) {
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
	if err = ctx.ShouldBindJSON(&questionInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	question := exam_question_domain.ExamQuestion{
		ID:            primitive.NewObjectID(),
		ExamID:        questionInput.ExamID,
		VocabularyID:  questionInput.VocabularyID,
		Content:       questionInput.Content,
		Type:          questionInput.Type,
		Level:         questionInput.Level,
		Options:       questionInput.Options,
		CorrectAnswer: questionInput.CorrectAnswer,
		CreatedAt:     time.Now(),
		UpdateAt:      time.Now(),
		WhoUpdate:     admin.FullName,
	}

	err = e.ExamQuestionUseCase.CreateOneInAdmin(ctx, &question)
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
