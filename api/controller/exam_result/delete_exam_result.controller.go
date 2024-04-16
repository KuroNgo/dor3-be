package exam_result_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExamResultController) DeleteOneExam(ctx *gin.Context) {
	answerID := ctx.Query("_id")

	err := e.ExamResultUseCase.DeleteOne(ctx, answerID)
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
