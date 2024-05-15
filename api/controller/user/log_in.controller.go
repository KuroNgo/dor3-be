package user_controller

import (
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (l *LoginFromRoleController) LoginFromRole(ctx *gin.Context) {
	//  Lấy thông tin từ request
	var adminInput admin_domain.SignIn
	if err := ctx.ShouldBindJSON(&adminInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	accessTokenCh := make(chan string)
	refreshTokenCh := make(chan string)

	var userInput user_domain.SignIn
	userInput.Email = adminInput.Email
	userInput.Password = adminInput.Password

	// Kiểm tra thông tin đăng nhập trong cả hai bảng user và admin
	admin, err := l.AdminUseCase.Login(ctx, adminInput)
	if err == nil && admin.Role == "admin" {
		go func() {
			defer close(accessTokenCh)
			// Generate token
			accessToken, err := internal.CreateToken(l.Database.AccessTokenExpiresIn, admin.Id, l.Database.AccessTokenPrivateKey)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "fail",
					"message": err.Error()},
				)
				return
			}
			accessTokenCh <- accessToken
		}()

		go func() {
			defer close(refreshTokenCh)
			refreshToken, err := internal.CreateToken(l.Database.RefreshTokenExpiresIn, admin.Id, l.Database.RefreshTokenPrivateKey)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "fail",
					"message": err.Error()},
				)
				return
			}
			refreshTokenCh <- refreshToken
		}()

		accessToken := <-accessTokenCh
		refreshToken := <-refreshTokenCh

		ctx.SetCookie("access_token", accessToken, l.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
		ctx.SetCookie("refresh_token", refreshToken, l.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
		ctx.SetCookie("logged_in", "true", l.Database.AccessTokenMaxAge*1000, "/", "localhost", false, false)

		ctx.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"message":      "Login successful with admin role",
			"access_token": accessToken,
		})
		return
	}

	// Tìm kiếm user trong database
	user, err := l.UserUseCase.Login(ctx, userInput)
	if err == nil && user.Verified == true {
		go func() {
			defer close(accessTokenCh)
			// Generate token
			accessToken, err := internal.CreateToken(l.Database.AccessTokenExpiresIn, user.ID, l.Database.AccessTokenPrivateKey)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "fail",
					"message": err.Error()},
				)
				return
			}
			accessTokenCh <- accessToken
		}()

		go func() {
			defer close(refreshTokenCh)
			refreshToken, err := internal.CreateToken(l.Database.RefreshTokenExpiresIn, user.ID, l.Database.RefreshTokenPrivateKey)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "fail",
					"message": err.Error()},
				)
				return
			}
			refreshTokenCh <- refreshToken
		}()

		accessToken := <-accessTokenCh
		refreshToken := <-refreshTokenCh

		ctx.SetCookie("access_token", accessToken, l.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
		ctx.SetCookie("refresh_token", refreshToken, l.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
		ctx.SetCookie("logged_in", "true", l.Database.AccessTokenMaxAge*1000, "/", "localhost", false, false)

		ctx.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"message":      "Login successful with user role",
			"access_token": accessToken,
		})
		return
	}

	// Trả về thông báo login không thành công
	ctx.JSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"message": err.Error(),
	})
}

func (l *LoginFromRoleController) LoginUser(ctx *gin.Context) {
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

	// Tìm kiếm user trong database
	user, err := l.UserUseCase.Login(ctx, userInput)
	if err == nil && user.Verified == true {
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

		ctx.SetCookie("access_token", accessToken, l.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
		ctx.SetCookie("refresh_token", refreshToken, l.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
		ctx.SetCookie("logged_in", "true", l.Database.AccessTokenMaxAge*1000, "/", "localhost", false, false)

		ctx.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"message":      "Login successful with user role",
			"access_token": accessToken,
		})
		return
	}

	// Trả về thông báo login không thành công
	ctx.JSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"message": err.Error(),
	})
}

func (a *LoginFromRoleController) LoginAdmin(ctx *gin.Context) {
	//  Lấy thông tin từ request
	var adminInput admin_domain.SignIn
	if err := ctx.ShouldBindJSON(&adminInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	// Kiểm tra thông tin đăng nhập trong cả hai bảng user và admin
	admin, err := a.AdminUseCase.Login(ctx, adminInput)
	if err == nil && admin.Role == "admin" {

		accessTokenCh := make(chan string)
		// Generate token
		go func() {
			defer close(accessTokenCh)
			accessToken, err := internal.CreateToken(a.Database.AccessTokenExpiresIn, admin.Id, a.Database.AccessTokenPrivateKey)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "fail",
					"message": err.Error()},
				)
				return
			}
			accessTokenCh <- accessToken
		}()

		refreshTokenCh := make(chan string)
		go func() {
			defer close(refreshTokenCh)
			refreshToken, err := internal.CreateToken(a.Database.RefreshTokenExpiresIn, admin.Id, a.Database.RefreshTokenPrivateKey)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "fail",
					"message": err.Error()},
				)
				return
			}
			refreshTokenCh <- refreshToken
		}()

		accessToken := <-accessTokenCh
		refreshToken := <-refreshTokenCh

		ctx.SetCookie("access_token", accessToken, a.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
		ctx.SetCookie("refresh_token", refreshToken, a.Database.AccessTokenMaxAge*1000, "/", "localhost", false, true)
		ctx.SetCookie("logged_in", "true", a.Database.AccessTokenMaxAge*1000, "/", "localhost", false, false)

		ctx.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"message":      "Login successful with admin role",
			"access_token": accessToken,
		})
		return
	}

	// Trả về thông báo login không thành công
	ctx.JSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"message": err.Error(),
	})
}
