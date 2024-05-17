package quiz_question_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (q *QuizQuestionsController) FetchManyQuizQuestion(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	quizID := ctx.Query("quiz_id")
	exam, err := q.QuizQuestionUseCase.FetchManyByQuizID(ctx, quizID)
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

func (q *QuizQuestionsController) FetchManyQuizQuestionInAdmin(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	quizID := ctx.Query("quiz_id")
	exam, err := q.QuizQuestionUseCase.FetchManyByQuizID(ctx, quizID)
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
func (q *QuizQuestionsController) FetchOneQuizQuestionByID(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	user, err := q.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	id := ctx.Query("_id")
	quiz, err := q.QuizQuestionUseCase.FetchByID(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   quiz,
	})
}
