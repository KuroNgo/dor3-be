package activity_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *ActivityController) FetchManyActivity(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	activity, err := a.ActivityUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"status":       "success",
		"activity_log": activity,
	})
}
