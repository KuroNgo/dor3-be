package exam_answer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExamAnswerController) DeleteOneExam(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	user, err := e.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	answerID := ctx.Query("_id")
	err = e.ExamAnswerUseCase.DeleteOne(ctx, answerID)
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
