package unit_controller

import (
	unit_domain "clean-architecture/domain/unit"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Deprecated: UpsertOneUnit
func (u *UnitController) UpsertOneUnit(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")

	user, err := u.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	unitID := ctx.Query("_id")

	var unitInput unit_domain.Unit
	if err := ctx.ShouldBindJSON(&unitInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	upsertUnit := unit_domain.Unit{
		LessonID:   unitInput.LessonID,
		Name:       unitInput.Name,
		Content:    unitInput.Content,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
		WhoUpdates: user.FullName,
	}

	unitRes, err := u.UnitUseCase.UpsertOne(ctx, unitID, &upsertUnit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   unitRes,
	})
}
