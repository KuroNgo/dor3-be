package admin_controller

import (
	"clean-architecture/internal"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *AdminController) RefreshToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": message,
		})
		return
	}

	sub, err := internal.ValidateToken(cookie, a.Database.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	admin, err := a.AdminUseCase.GetByID(ctx, fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "the user belonging to this token no logger exists",
		})
		return
	}

	access_token, err := internal.CreateToken(a.Database.AccessTokenExpiresIn, admin.Id, a.Database.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	ctx.SetCookie("access_token", access_token, a.Database.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", a.Database.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"access_token": access_token,
	})
}
