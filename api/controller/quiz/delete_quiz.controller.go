package quiz_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOneQuiz done
func (q *QuizController) DeleteOneQuiz(ctx *gin.Context) {
	quizID := ctx.Query("_id")

	err := q.QuizUseCase.DeleteOne(ctx, quizID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Trả về mảng dữ liệu dưới dạng JSON
	ctx.JSON(http.StatusOK, "the quiz is deleted!")

}
