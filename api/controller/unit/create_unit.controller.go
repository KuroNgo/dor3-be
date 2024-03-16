package unit_controller

import (
	unit_domain "clean-architecture/domain/unit"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (u *UnitController) CreateOneUnit(ctx *gin.Context) {
	var unitInput unit_domain.Input
	if err := ctx.ShouldBindJSON(&unitInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := internal.IsValidUnit(unitInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	unitRes := &unit_domain.Unit{
		ID:        primitive.NewObjectID(),
		LessonID:  unitInput.LessonID,
		Name:      unitInput.Name,
		Content:   unitInput.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		//WhoUpdates:
	}

	err := u.UnitUseCase.CreateOne(ctx, unitRes)
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
