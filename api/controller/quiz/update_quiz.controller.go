package quiz_controller

import (
	quiz_domain "clean-architecture/domain/quiz"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateOneQuiz done
func (q *QuizController) UpdateOneQuiz(ctx *gin.Context) {
	quizID := ctx.Query("_id")

	var quiz quiz_domain.Input
	if err := ctx.ShouldBindJSON(&quiz); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	updateQuiz := quiz_domain.Quiz{
		Question:      quiz.Question,
		Options:       quiz.Options,
		CorrectAnswer: quiz.CorrectAnswer,
		Explanation:   quiz.Explanation,
		QuestionType:  quiz.QuestionType,
	}

	err := q.QuizUseCase.UpdateOne(ctx, quizID, updateQuiz)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
