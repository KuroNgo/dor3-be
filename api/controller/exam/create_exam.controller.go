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

	lessonID := ctx.Query("lesson_id")
	idLesson, _ := primitive.ObjectIDFromHex(lessonID)

	UnitID := ctx.Query("unit_id")
	idUnit, _ := primitive.ObjectIDFromHex(UnitID)

	VocabularyID := ctx.Query("vocabulary_id")
	idVocabulary, _ := primitive.ObjectIDFromHex(VocabularyID)

	var examInput exam_domain.Input
	if err := ctx.ShouldBindJSON(&examInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	exam := exam_domain.Exam{
		ID:           primitive.NewObjectID(),
		LessonID:     idLesson,
		UnitID:       idUnit,
		VocabularyID: idVocabulary,

		Title:       examInput.Title,
		Description: examInput.Description,
		Duration:    examInput.Duration,

		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		WhoUpdates: user.FullName,
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
