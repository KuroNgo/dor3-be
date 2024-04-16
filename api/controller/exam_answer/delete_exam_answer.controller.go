package exam_answer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExamAnswerController) DeleteOneExam(ctx *gin.Context) {
	answerID := ctx.Query("_id")

	err := e.ExamAnswerUseCase.DeleteOne(ctx, answerID)
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
