package lesson_controller

import (
	lesson_domain "clean-architecture/domain/lesson"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (l *LessonController) UpdateOneLesson(ctx *gin.Context) {
	lessonID := ctx.Query("_id")

	var lessonInput lesson_domain.Input
	if err := ctx.ShouldBindJSON(&lessonInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	updateLesson := lesson_domain.Lesson{
		CourseID:  lessonInput.CourseID,
		Name:      lessonInput.Name,
		Content:   lessonInput.Content,
		UpdatedAt: time.Now(),
		//WhoUpdates:
	}

	err := l.LessonUseCase.UpdateOne(ctx, lessonID, updateLesson)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}