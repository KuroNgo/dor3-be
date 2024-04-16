package exam_answer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExamAnswerController) FetchManyExam(ctx *gin.Context) {
	questionID := ctx.Query("question_id")

	answer, err := e.ExamAnswerUseCase.FetchManyByQuestionID(ctx, questionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   answer,
	})
}
