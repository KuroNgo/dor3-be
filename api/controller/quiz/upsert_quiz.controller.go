package quiz_controller

import (
	quiz_domain "clean-architecture/domain/quiz"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Deprecated: UpsertOneQuiz
func (q *QuizController) UpsertOneQuiz(ctx *gin.Context) {
	// Kiểm tra xác thực người dùng
	currentUser := ctx.MustGet("currentUser")
	user, err := q.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting user: " + err.Error(),
		})
		return
	}

	// Lấy ID của quiz (nếu có)
	quizID := ctx.Query("_id")

	// Bind dữ liệu JSON từ yêu cầu HTTP vào biến quiz
	var quiz quiz_domain.Input
	if err := ctx.ShouldBindJSON(&quiz); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Error binding JSON: " + err.Error(),
		})
		return
	}

	// Tạo hoặc cập nhật bài kiểm tra
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
	if quizID != "" {
		quizRes, err = q.QuizUseCase.UpsertOne(ctx, quizID, &upsertQuiz)
	} else {
		err = q.QuizUseCase.CreateOne(ctx, &upsertQuiz)
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Error binding JSON: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   quizRes,
	})
}
