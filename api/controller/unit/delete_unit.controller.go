package unit_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *UnitController) DeleteOneUnit(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := u.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	lessonID := ctx.Query("_id")

	err = u.UnitUseCase.DeleteOne(ctx, lessonID)
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
