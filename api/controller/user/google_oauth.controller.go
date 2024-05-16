package user_controller

import (
	"clean-architecture/bootstrap"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal"
	google_utils "clean-architecture/internal/Oauth2/google"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"time"
)

type GoogleAuthController struct {
	GoogleAuthUseCase user_domain.IGoogleAuthUseCase
	Database          *bootstrap.Database
}

func (auth *GoogleAuthController) GoogleLoginWithUser(c *gin.Context) {
	var googleOauthConfig = &oauth2.Config{}
	googleOauthConfig = &oauth2.Config{
		ClientID:     auth.Database.GoogleClientID,
		ClientSecret: auth.Database.GoogleClientSecret,
		RedirectURL:  auth.Database.GoogleOAuthRedirectUrl,
		Scopes:       []string{"profile", "email"}, // Adjust scopes as needed
		Endpoint:     google.Endpoint,
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Error exchanging code: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInfo, err := google_utils.GetUserInfo(token.AccessToken)
	if err != nil {
		fmt.Println("Error getting user info: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Giả sử userInfo là một map[string]interface{}
	email := userInfo["email"].(string)
	fullName := userInfo["name"].(string)
	avatarURL := userInfo["picture"].(string)
	verifiedEmail := userInfo["verified_email"].(bool)
	resBody := &user_domain.UserInput{
		ID:        primitive.NewObjectID(),
		Email:     email,
		FullName:  fullName,
		AvatarURL: avatarURL,
		Provider:  "google",
		Role:      "user",
		Verified:  verifiedEmail,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updateUser, err := auth.GoogleAuthUseCase.UpsertUser(c, resBody.Email, resBody)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	signedToken, err := google_utils.SignJWT(userInfo)
	if err != nil {
		fmt.Println("Error signing token: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := internal.CreateToken(auth.Database.AccessTokenExpiresIn, updateUser.ID.Hex(), auth.Database.AccessTokenPrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	refreshToken, err := internal.CreateToken(auth.Database.RefreshTokenExpiresIn, updateUser.ID.Hex(), auth.Database.RefreshTokenPrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", accessToken, auth.Database.AccessTokenMaxAge*60, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, auth.Database.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	c.SetCookie("logged_in", "true", auth.Database.AccessTokenMaxAge*60, "/", "localhost", false, false)
	c.SetSameSite(http.SameSiteStrictMode) // Chỉ gửi cookie với các yêu cầu cùng nguồn

	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}
