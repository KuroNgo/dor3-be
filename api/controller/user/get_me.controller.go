package user_controller

import (
	"clean-architecture/internal"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *UserController) GetMe(ctx *gin.Context) {
	cookie, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	sub, err := internal.ValidateToken(cookie, u.Database.AccessTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	result, err := u.UserUseCase.GetByID(ctx, fmt.Sprint(sub))
	resultString, err := json.Marshal(result)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": string(resultString) + "the user belonging to this token no logger exists",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   result,
	})
}
