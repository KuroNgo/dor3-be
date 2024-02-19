package quiz_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (q *QuizController) DeleteQuiz(ctx *gin.Context) {
	quizID := ctx.Param("_id")

	err := q.QuizUseCase.Delete(ctx, quizID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Trả về mảng dữ liệu dưới dạng JSON
	ctx.JSON(http.StatusOK, "the quiz is deleted!")

}
