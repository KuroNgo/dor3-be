package unit_controller

import (
	unit_domain "clean-architecture/domain/_unit"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (u *UnitController) UpdateOneUnit(ctx *gin.Context) {
	lessonID := ctx.Query("_id")

	var unitInput unit_domain.Input
	if err := ctx.ShouldBindJSON(&unitInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	updateLesson := unit_domain.Unit{
		LessonID:  unitInput.LessonID,
		Name:      unitInput.Name,
		Content:   unitInput.Content,
		UpdatedAt: time.Now(),
		//WhoUpdates:
	}

	err := u.UnitUseCase.UpdateOne(ctx, lessonID, updateLesson)
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
