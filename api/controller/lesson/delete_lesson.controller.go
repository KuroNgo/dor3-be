package lesson_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOneLesson chỉ admin mới có thể xóa
func (l *LessonController) DeleteOneLesson(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := l.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": admin.FullName + "You are not authorized to perform this action!",
		})
		return
	}

	lessonID := ctx.Param("_id")

	err = l.LessonUseCase.DeleteOne(ctx, lessonID)
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
