package quiz_answer_controller

import (
	quiz_answer_domain "clean-architecture/domain/quiz_answer"
	quiz_result_domain "clean-architecture/domain/quiz_result"
	user_attempt_domain "clean-architecture/domain/user_process/exam_management"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

// handleError xử lý lỗi chung
func handleError(ctx *gin.Context, statusCode int, message string, err error) {
	ctx.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
		"error":   err.Error(),
	})
}

// IsCorrect là hằng số để biểu thị câu trả lời đúng
const IsCorrect = 1

func (q *QuizAnswerController) CreateOneQuizAnswer(ctx *gin.Context) {
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

	var answerInput quiz_answer_domain.Input
	if err = ctx.ShouldBindJSON(&answerInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	answer := quiz_answer_domain.QuizAnswer{
		ID:          primitive.NewObjectID(),
		UserID:      user.ID,
		QuestionID:  answerInput.QuestionID,
		Answer:      answerInput.Answer,
		SubmittedAt: time.Now(),
	}

	err = q.QuizAnswerUseCase.CreateOne(ctx, &answer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	data, err := q.QuizAnswerUseCase.FetchManyAnswerByUserIDAndQuestionID(ctx, answerInput.QuestionID.Hex(), user.ID.Hex())
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Failed to fetch answers", err)
	}

	var quizID primitive.ObjectID
	var totalCorrect int16

	for i, res := range data.QuizAnswer {
		if i == 0 {
			quizID = res.Question.QuizID
		}
		if res.IsCorrect == IsCorrect {
			totalCorrect++
		}
	}

	if quizID != primitive.NilObjectID {
		quizResult := &quiz_result_domain.QuizResult{
			ID:         primitive.NewObjectID(),
			UserID:     user.ID,
			QuizID:     quizID,
			Score:      totalCorrect,
			StartedAt:  time.Now(),
			IsComplete: IsCorrect,
		}

		err := q.QuizResultUseCase.CreateOne(ctx, quizResult)
		if err != nil {
			handleError(ctx, http.StatusInternalServerError, "Failed to create exam result", err)
			return
		}

		userProcess := user_attempt_domain.ExamManagement{
			ID:            primitive.NewObjectID(),
			UserID:        user.ID,
			QuizID:        quizID,
			Score:         float32(totalCorrect),
			ProcessStatus: 0,
			CompletedDate: time.Now(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		err = q.UserAttemptUseCase.UpdateExamManagementByQuizID(ctx, userProcess)
		if err != nil {
			handleError(ctx, http.StatusInternalServerError, "Failed to create user process", err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
