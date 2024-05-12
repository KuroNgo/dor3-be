package unit_controller

import (
	unit_domain "clean-architecture/domain/unit"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (u *UnitController) UpdateCompleteUnit(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	user, err := u.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var unitInput unit_domain.CompleteInput
	if err := ctx.ShouldBindJSON(&unitInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	updateUnit := &unit_domain.Unit{
		ID:         unitInput.ID,
		IsComplete: unitInput.IsComplete,
		UpdatedAt:  time.Now(),
		Learner:    user.FullName,
	}

	err = u.UnitUseCase.UpdateComplete(ctx, updateUnit)
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
