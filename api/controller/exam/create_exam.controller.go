package exam_controller

import (
	exam_domain "clean-architecture/domain/exam"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (e *ExamsController) CreateOneExam(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")

	user, err := e.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	var examInput exam_domain.Input
	if err := ctx.ShouldBindJSON(&examInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	exam := exam_domain.Exam{
		ID:            primitive.NewObjectID(),
		LessonID:      examInput.LessonID,
		UnitID:        examInput.UnitID,
		VocabularyID:  examInput.VocabularyID,
		Question:      examInput.Question,
		Options:       examInput.Options,
		CorrectAnswer: examInput.CorrectAnswer,
		Explanation:   examInput.Explanation,
		QuestionType:  examInput.QuestionType,
		Level:         examInput.Level,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		WhoUpdates:    user.FullName,
	}

	err = e.ExamUseCase.CreateOne(ctx, &exam)
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
