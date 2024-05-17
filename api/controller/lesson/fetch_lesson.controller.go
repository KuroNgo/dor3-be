package lesson_controller

import (
	"fmt"
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

func (l *LessonController) FetchManyNotPagination(ctx *gin.Context) {
	lesson, err := l.LessonUseCase.FetchManyNotPagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
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

func (l *LessonController) FetchManyInAdmin(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := l.AdminUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

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
func (l *LessonController) FetchByIdCourseInAdmin(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := l.AdminUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

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

func (l *LessonController) FetchByIdInAdmin(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := l.AdminUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

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
