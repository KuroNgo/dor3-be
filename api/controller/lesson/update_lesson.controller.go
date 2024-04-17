package lesson_controller

import (
	lesson_domain "clean-architecture/domain/lesson"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (l *LessonController) UpdateOneLesson(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	user, err := l.AdminUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var lessonInput lesson_domain.Input
	if err := ctx.ShouldBindJSON(&lessonInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	updateLesson := lesson_domain.Lesson{
		ID:       lessonInput.ID,
		CourseID: lessonInput.CourseID,

		Name:       lessonInput.Name,
		Content:    lessonInput.Content,
		UpdatedAt:  time.Now(),
		WhoUpdates: user.FullName,
	}

	data, err := l.LessonUseCase.UpdateOne(ctx, &updateLesson)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"result": data,
	})
}
