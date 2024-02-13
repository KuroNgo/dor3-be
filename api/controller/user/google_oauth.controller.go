package user_controller

import (
	"clean-architecture/bootstrap"
	user_domain "clean-architecture/domain/request/user"
	"clean-architecture/internal"
	"clean-architecture/internal/Oauth2/google"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GoogleAuthController struct {
	GoogleAuthUseCase user_domain.IGoogleAuthUseCase
	Database          *bootstrap.Database
}

func (auth *GoogleAuthController) GoogleLogin(ctx *gin.Context) {
	code := ctx.Query("code")
	pathUrl := "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "Authorization code not provided!",
		})
		return
	}

	// Use the code get the id and access tokens
	tokenRes, err := google.GetGoogleOauthToken(code)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	user, err := google.GetGoogleUser(tokenRes.AccessToken, tokenRes.IDToken)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	createdAt := time.Now()
	resBody := &user_domain.User{
		Email:     user.Email,
		FullName:  user.Name,
		AvatarURL: user.Picture,
		Provider:  "google",
		Role:      "user",
		Verified:  true,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	updatedUser, err := auth.GoogleAuthUseCase.UpsertUser(ctx, user.Email, resBody)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	app := bootstrap.App()

	env := app.Env
	accessToken, err := internal.CreateToken(env.AccessTokenExpiresIn, updatedUser.ID.Hex(), env.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := internal.CreateToken(env.RefreshTokenExpiresIn, updatedUser.ID.Hex(), env.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, env.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, env.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", env.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(env.ClientOrigin, pathUrl))
}
