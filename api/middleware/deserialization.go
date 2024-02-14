package middleware

import (
	"clean-architecture/bootstrap"
	"clean-architecture/internal"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": "You are not logged in",
			})
			return
		}

		app := bootstrap.App()
		env := app.Env

		sub, err := internal.ValidateToken(accessToken, env.AccessTokenPublicKey)

		if err != nil {
			fmt.Println("The err is: ", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}

		//var user model.User
		//
		//result := conf.DbDefault.First(&user, "userid = ?", fmt.Sprint(sub))
		//if result.Error != nil {
		//	ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		//	return
		//}

		ctx.Set("currentUser", sub)
		ctx.Next()
	}
}
