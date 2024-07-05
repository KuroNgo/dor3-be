package exam_answer_controller

import (
	exam_answer_domain "clean-architecture/domain/exam_answer"
	exam_result_domain "clean-architecture/domain/exam_result"
	user_attempt_domain "clean-architecture/domain/user_process"
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

// CreateOneExamAnswer tạo một câu trả lời bài kiểm tra
func (e *ExamAnswerController) CreateOneExamAnswer(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		handleError(ctx, http.StatusUnauthorized, "You are not logged in!", nil)
		return
	}

	user, err := e.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		handleError(ctx, http.StatusUnauthorized, "You are not authorized to perform this action!", err)
		return
	}

	var answerInput exam_answer_domain.Input
	if err = ctx.ShouldBindJSON(&answerInput); err != nil {
		handleError(ctx, http.StatusBadRequest, "Invalid request", err)
		return
	}

	answer := exam_answer_domain.ExamAnswer{
		ID:          primitive.NewObjectID(),
		UserID:      user.ID,
		QuestionID:  answerInput.QuestionID,
		Answer:      answerInput.Answer,
		SubmittedAt: time.Now(),
	}

	err = e.ExamAnswerUseCase.CreateOneInUser(ctx, &answer)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Failed to create answer", err)
		return
	}

	data, err := e.ExamAnswerUseCase.FetchManyAnswerByQuestionIDInUser(ctx, answerInput.QuestionID.Hex(), user.ID)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Failed to fetch answers", err)
	}

	var examID primitive.ObjectID
	var totalCorrect int16

	for i, res := range data.ExamAnswerResponse {
		if i == 0 {
			examID = res.Question.ExamID
		}
		if res.IsCorrect == IsCorrect {
			totalCorrect++
		}
	}

	if examID != primitive.NilObjectID {
		examResult := &exam_result_domain.ExamResult{
			ID:         primitive.NewObjectID(),
			UserID:     user.ID,
			ExamID:     examID,
			Score:      totalCorrect,
			StartedAt:  time.Now(),
			IsComplete: IsCorrect,
		}

		err := e.ExamResultUseCase.CreateOneInUser(ctx, examResult)
		if err != nil {
			handleError(ctx, http.StatusInternalServerError, "Failed to create exam result", err)
			return
		}

		userProcess := user_attempt_domain.ExamManagement{
			ID:            primitive.NewObjectID(),
			UserID:        user.ID,
			ExamID:        examID,
			QuizID:        primitive.NilObjectID,
			ExerciseID:    primitive.NilObjectID,
			Score:         float32(totalCorrect),
			ProcessStatus: 0,
			CompletedDate: time.Now(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		err = e.UserAttemptUseCase.UpdateExamManagementByExamID(ctx, userProcess)
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
