package unit_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *UnitController) DeleteOneUnit(ctx *gin.Context) {
	lessonID := ctx.Query("_id")

	err := u.UnitUseCase.DeleteOne(ctx, lessonID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Trả về mảng dữ liệu dưới dạng JSON
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "the unit is deleted!",
	})
}
