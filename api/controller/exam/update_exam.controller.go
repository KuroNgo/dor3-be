package exam_controller

import (
	exam_domain "clean-architecture/domain/exam"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (e *ExamsController) UpdateOneExam(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := e.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var examInput exam_domain.Input
	if err := ctx.ShouldBindJSON(&examInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	exam := exam_domain.Exam{
		ID:           examInput.ID,
		LessonID:     examInput.LessonID,
		UnitID:       examInput.UnitID,
		VocabularyID: examInput.VocabularyID,

		Title:       examInput.Title,
		Description: examInput.Description,
		Duration:    examInput.Duration,

		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		WhoUpdates: admin.FullName,
	}

	_, err = e.ExamUseCase.UpdateOne(ctx, &exam)
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
