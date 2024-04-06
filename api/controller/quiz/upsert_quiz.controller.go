package quiz_controller

import (
	quiz_domain "clean-architecture/domain/quiz"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (q *QuizController) UpsertOneQuiz(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	user, err := q.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user: " + err.Error()})
		return
	}

	quizID := ctx.Query("_id")

	var quiz quiz_domain.Input
	if err := ctx.ShouldBindJSON(&quiz); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Error binding JSON: " + err.Error(),
		})
		return
	}

	upsertQuiz := quiz_domain.Quiz{
		Question:      quiz.Question,
		Options:       quiz.Options,
		CorrectAnswer: quiz.CorrectAnswer,
		Explanation:   quiz.Explanation,
		QuestionType:  quiz.QuestionType,
		UpdatedAt:     time.Now(),
		WhoUpdates:    user.FullName,
	}

	var quizRes quiz_domain.Response
	var upsertErr error
	if quizID != "" {
		quizRes, upsertErr = q.QuizUseCase.UpsertOne(ctx, quizID, &upsertQuiz)
	} else {
		upsertErr = q.QuizUseCase.CreateOne(ctx, &upsertQuiz)
	}

	if upsertErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error upserting quiz: " + upsertErr.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": quizRes})
}
