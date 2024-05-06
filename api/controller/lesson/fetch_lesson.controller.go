package lesson_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (l *LessonController) FetchMany(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	lesson, detail, err := l.LessonUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"detail": detail,
		"data":   lesson,
	})
}

func (l *LessonController) FetchByIdCourse(ctx *gin.Context) {
	idCourse := ctx.Query("course_id")
	page := ctx.DefaultQuery("page", "1")

	lesson, detail, err := l.LessonUseCase.FetchByIdCourse(ctx, idCourse, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"detail": detail,
		"lesson": lesson,
	})
}

func (l *LessonController) FetchById(ctx *gin.Context) {
	idLesson := ctx.Query("_id")
	lesson, err := l.LessonUseCase.FetchByID(ctx, idLesson)
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
