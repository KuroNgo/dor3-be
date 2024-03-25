package lesson_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (l *LessonController) FetchMany(ctx *gin.Context) {
	lesson, err := l.LessonUseCase.FetchMany(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"lesson": lesson,
		},
	})
}

func (l *LessonController) FetchByIdCourse(ctx *gin.Context) {
	idCourse := ctx.Param("course_id")
	lesson, err := l.LessonUseCase.FetchByIdCourse(ctx, idCourse)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"lesson": lesson,
	})
}
