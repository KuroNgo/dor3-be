package admin_controller

import (
	"clean-architecture/internal"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *AdminController) GetMe(ctx *gin.Context) {
	cookie, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	sub, err := internal.ValidateToken(cookie, a.Database.AccessTokenPublicKey)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	result, err := a.AdminUseCase.GetByID(ctx, fmt.Sprint(sub))
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "Failed to get user data: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   result,
	})
}
