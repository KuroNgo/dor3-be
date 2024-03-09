package lesson_controller

import (
	lesson_domain "clean-architecture/domain/lesson"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (l *LessonController) UpsertOneLesson(ctx *gin.Context) {
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
		CourseID:  lessonInput.CourseID,
		Name:      lessonInput.Name,
		Content:   lessonInput.Content,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		//WhoUpdates:
	}

	lessRes, err := l.LessonUseCase.UpsertOne(ctx, lessonID, &upsertLesson)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   lessRes,
	})
}
