package quiz_options_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (q *QuizOptionsController) FetchManyExerciseOptions(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	questionID := ctx.Query("question_id")
	exam, err := q.QuizOptionsUseCase.FetchManyByQuestionID(ctx, questionID)
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
