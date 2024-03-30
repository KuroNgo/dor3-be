package user_controller

import (
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (l *LoginFromRoleController) Login2(ctx *gin.Context) {
	//  Lấy thông tin từ request
	var adminInput admin_domain.SignIn
	if err := ctx.ShouldBindJSON(&adminInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	var userInput user_domain.SignIn
	userInput.Email = adminInput.Email
	userInput.Password = adminInput.Password

	// Kiểm tra thông tin đăng nhập trong cả hai bảng user và admin
	admin, err := l.AdminUseCase.Login(ctx, adminInput)
	if err == nil && admin.Role == "admin" {
		// Generate token
		accessToken, err := internal.CreateToken(l.Database.AccessTokenExpiresIn, admin.Id, l.Database.AccessTokenPrivateKey)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": err.Error()},
			)
			return
		}

		refreshToken, err := internal.CreateToken(l.Database.RefreshTokenExpiresIn, admin.Id, l.Database.RefreshTokenPrivateKey)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": err.Error()},
			)
			return
		}

		ctx.SetCookie("access_token", accessToken, l.Database.AccessTokenMaxAge*60, "/", "localhost", false, true)
		ctx.SetCookie("refresh_token", refreshToken, l.Database.RefreshTokenMaxAge*60, "/", "localhost", false, true)
		ctx.SetCookie("logged_in", "true", l.Database.AccessTokenMaxAge*60, "/", "localhost", false, true)

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Login successful with admin role",
		})
		return
	}

	// Tìm kiếm user trong database
	user, err := l.UserUseCase.Login(ctx, userInput)
	if err == nil && user.Role == "user" {
		// Generate token
		accessToken, err := internal.CreateToken(l.Database.AccessTokenExpiresIn, user.ID, l.Database.AccessTokenPrivateKey)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": err.Error()},
			)
			return
		}

		refreshToken, err := internal.CreateToken(l.Database.RefreshTokenExpiresIn, user.ID, l.Database.RefreshTokenPrivateKey)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": err.Error()},
			)
			return
		}

		ctx.SetCookie("access_token", accessToken, l.Database.AccessTokenMaxAge*60, "/", "localhost", false, true)
		ctx.SetCookie("refresh_token", refreshToken, l.Database.RefreshTokenMaxAge*60, "/", "localhost", false, true)
		ctx.SetCookie("logged_in", "true", l.Database.AccessTokenMaxAge*60, "/", "localhost", false, true)

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Login successful with user role",
		})
		return
	}

	// Trả về thông báo login không thành công
	ctx.JSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"message": err.Error(),
	})
}
