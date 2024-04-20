package exam_question_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExamQuestionsController) FetchManyExamOptions(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	examID := ctx.Query("exam_id")
	exam, err := e.ExamQuestionUseCase.FetchManyByExamID(ctx, examID)
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
