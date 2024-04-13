package exam_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExamsController) DeleteOneExam(ctx *gin.Context) {
	examID := ctx.Query("_id")

	err := e.ExamUseCase.DeleteOne(ctx, examID)
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
