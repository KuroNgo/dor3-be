package lesson_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOneLesson chỉ admin mới có thể xóa
func (l *LessonController) DeleteOneLesson(ctx *gin.Context) {
	lessonID := ctx.Query("_id")

	err := l.LessonUseCase.DeleteOne(ctx, lessonID)
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
