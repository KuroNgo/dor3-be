package exam_result_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExamResultController) FetchManyExam(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	exam, err := e.ExamResultUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   exam,
	})
}

func (e *ExamResultController) FetchResultByExamID(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	exam, err := e.ExamResultUseCase.FetchManyByExamID(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   exam,
	})
}

func (e *ExamResultController) GetResultsByUserIDAndExamID(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	examId := ctx.Query("exam_id")

	exam, err := e.ExamResultUseCase.GetResultsByUserIDAndExamID(ctx, fmt.Sprintf("%s", currentUser), examId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   exam,
	})
}
