package exam_result_controller

import (
	user_attempt_domain "clean-architecture/domain/user_process"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

func (e *ExamResultController) DeleteOneExamResult(ctx *gin.Context) {
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

	resultID := ctx.Query("_id")
	if resultID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Result ID is required",
		})
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = e.ExamResultUseCase.DeleteOneInUser(ctx, resultID)
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
		res, err := e.ExamResultUseCase.GetResultsByExamIDInUser(ctx, user.ID.Hex(), resultID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		userAttempt := user_attempt_domain.ExamManagement{
			UserID:        res.UserID,
			ExamID:        res.ExamID,
			Score:         0,
			ProcessStatus: 0,
			CompletedDate: time.Now(),
			UpdatedAt:     time.Now(),
		}
		err = e.UserAttemptUseCase.UpdateExamManagementByUserID(ctx, userAttempt)
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
