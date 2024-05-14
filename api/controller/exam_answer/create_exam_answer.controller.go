package exam_answer_controller

import (
	exam_answer_domain "clean-architecture/domain/exam_answer"
	exam_result_domain "clean-architecture/domain/exam_result"
	user_attempt_domain "clean-architecture/domain/user_attempt"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (e *ExamAnswerController) CreateOneExamAnswer(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	user, err := e.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var answerInput exam_answer_domain.Input
	if err = ctx.ShouldBindJSON(&answerInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	answer := exam_answer_domain.ExamAnswer{
		ID:          primitive.NewObjectID(),
		UserID:      user.ID,
		QuestionID:  answerInput.QuestionID,
		Answer:      answerInput.Answer,
		SubmittedAt: time.Now(),
	}

	err = e.ExamAnswerUseCase.CreateOne(ctx, &answer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	data, err := e.ExamAnswerUseCase.FetchManyAnswerByUserIDAndQuestionID(ctx, answerInput.QuestionID.Hex(), user.ID.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if len(data.ExamAnswerResponse) == 9 {
		var examID primitive.ObjectID
		var totalCorrect int16

		// Determine the examID and count the total correct answers
		for i, res := range data.ExamAnswerResponse {
			if i == 1 {
				examID = res.Question.ExamID
			}
			if res.IsCorrect == 1 {
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
				IsComplete: 1,
			}

			err := e.ExamResultUseCase.CreateOne(ctx, examResult)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}

			userProcess := user_attempt_domain.UserProcess{
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

			err = e.UserAttemptUseCase.CreateOneByUserID(ctx, userProcess)
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
	}
}
