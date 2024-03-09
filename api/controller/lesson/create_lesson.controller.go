package lesson_controller

import (
	lesson_domain "clean-architecture/domain/lesson"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (l *LessonController) CreateOneLesson(ctx *gin.Context) {
	var lessonInput lesson_domain.Input
	if err := ctx.ShouldBindJSON(&lessonInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if err := internal.IsValidLesson(lessonInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	lessonRes := &lesson_domain.Lesson{
		ID:        primitive.NewObjectID(),
		CourseID:  lessonInput.CourseID,
		Name:      lessonInput.Name,
		Content:   lessonInput.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		//WhoUpdates:
	}

	err := l.LessonUseCase.CreateOne(ctx, lessonRes)
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
