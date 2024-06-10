package exam_result_controller

import (
	exam_result_domain "clean-architecture/domain/exam_result"
	user_attempt_domain "clean-architecture/domain/user_process"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"sync"
	"time"
)

func (e *ExamResultController) CreateOneExamResult(ctx *gin.Context) {
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

	examID := ctx.Query("exam_id")
	idExam, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%s", examID))

	var inputResult exam_result_domain.Input
	if err := ctx.ShouldBindJSON(&inputResult); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var wg sync.WaitGroup
	var mutex sync.RWMutex
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := exam_result_domain.ExamResult{
			ID:         primitive.NewObjectID(),
			UserID:     user.ID,
			ExamID:     idExam,
			Score:      inputResult.Score,
			StartedAt:  inputResult.StartedAt,
			IsComplete: 1,
		}

		mutex.Lock()
		err = e.ExamResultUseCase.CreateOne(ctx, &result)
		mutex.Unlock()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		userAttempt := user_attempt_domain.ExamManagement{
			UserID:        user.ID,
			ExamID:        idExam,
			Score:         float32(inputResult.Score / 2),
			ProcessStatus: 1,
			CompletedDate: time.Now(),
			UpdatedAt:     time.Now(),
		}

		mutex.Lock()
		err = e.UserAttemptUseCase.UpdateExamManagementByUserID(ctx, userAttempt)
		mutex.Unlock()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	}()
	wg.Wait()

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
