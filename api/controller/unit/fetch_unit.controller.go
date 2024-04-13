package unit_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *UnitController) FetchMany(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	unit, err := u.UnitUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"unit": unit,
		},
	})
}

func (u *UnitController) FetchByIdLesson(ctx *gin.Context) {
	idLesson := ctx.Query("lesson_id")
	unit, err := u.UnitUseCase.FetchByIdLesson(ctx, idLesson)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"unit":   unit,
	})
}
