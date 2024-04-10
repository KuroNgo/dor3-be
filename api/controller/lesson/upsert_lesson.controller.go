package lesson_controller

import (
	lesson_domain "clean-architecture/domain/lesson"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Deprecated: UpsertOneLesson
func (l *LessonController) UpsertOneLesson(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")

	user, err := l.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	lessonID := ctx.Query("_id")

	var lessonInput lesson_domain.Lesson
	if err := ctx.ShouldBindJSON(&lessonInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	upsertLesson := lesson_domain.Lesson{
		CourseID:   lessonInput.CourseID,
		Name:       lessonInput.Name,
		Content:    lessonInput.Content,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
		WhoUpdates: user.FullName,
	}

	lessonRes, err := l.LessonUseCase.UpsertOne(ctx, lessonID, &upsertLesson)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   lessonRes,
	})
}
