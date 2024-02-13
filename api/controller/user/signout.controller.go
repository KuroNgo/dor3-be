package user_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
