package user_controller

import (
	"clean-architecture/bootstrap"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type GoogleAuthController struct {
	GoogleAuthUseCase user_domain.IGoogleAuthUseCase
	Database          *bootstrap.Database
}

func (auth *GoogleAuthController) GoogleLoginWithUser(ctx *gin.Context) {

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
	tokenRes, err := auth.GoogleAuthUseCase.GetGoogleOauthToken(code)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	user, err := auth.GoogleAuthUseCase.GetGoogleUser(tokenRes.AccessToken, tokenRes.IDToken)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	var role string
	if strings.Contains(user.Email, "feit") {
		role = "admin"
	} else {
		role = "user"
	}

	createdAt := time.Now()
	resBody := &user_domain.User{
		Email:     user.Email,
		FullName:  user.Name,
		AvatarURL: user.Picture,
		Provider:  "google",
		Role:      role,
		Verified:  true,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	updatedUser, err := auth.GoogleAuthUseCase.UpsertUser(ctx, user.Email, resBody)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	accessToken, err := internal.CreateToken(auth.Database.AccessTokenExpiresIn, updatedUser.ID.Hex(), auth.Database.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	refreshToken, err := internal.CreateToken(auth.Database.RefreshTokenExpiresIn, updatedUser.ID.Hex(), auth.Database.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	ctx.SetCookie("access_token", accessToken, auth.Database.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, auth.Database.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", auth.Database.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(auth.Database.ClientOrigin, pathUrl))
}
