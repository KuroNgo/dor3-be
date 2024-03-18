package unit_controller

import (
	unit_domain "clean-architecture/domain/unit"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (u *UnitController) UpsertOneUnit(ctx *gin.Context) {
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
		LessonID:  unitInput.LessonID,
		Name:      unitInput.Name,
		Content:   unitInput.Content,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		//WhoUpdates:
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
