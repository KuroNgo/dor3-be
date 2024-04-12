package admin_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *AdminController) Logout(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
