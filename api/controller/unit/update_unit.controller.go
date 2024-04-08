package unit_controller

import (
	unit_domain "clean-architecture/domain/unit"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (u *UnitController) UpdateOneUnit(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")

	user, err := u.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	unitID := ctx.Param("_id")

	var unitInput unit_domain.Input
	if err := ctx.ShouldBindJSON(&unitInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	updateUnit := unit_domain.Unit{
		LessonID:   unitInput.LessonID,
		Name:       unitInput.Name,
		Content:    unitInput.Content,
		UpdatedAt:  time.Now(),
		WhoUpdates: user.FullName,
	}

	unit, err := u.UnitUseCase.UpdateOne(ctx, unitID, updateUnit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   unit,
	})
}

func (u *UnitController) UpdateCompletedByIdUnit(ctx *gin.Context) {
	currentUser := ctx.MustGet("access_token")
	unitID := ctx.Param("_id")
	lessonID := ctx.Param("lesson_id")

	isComplete := 1

	user, err := u.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))

	updateReq := unit_domain.Update{
		UnitID:     unitID,
		LessonID:   lessonID,
		IsComplete: isComplete,
		WhoUpdate:  user.FullName,
	}
	err = u.UnitUseCase.UpdateComplete(ctx, updateReq)
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
