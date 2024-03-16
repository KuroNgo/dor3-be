package user_controller

import (
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *UserController) LogIn(ctx *gin.Context) {

}

func (l *LoginFromRoleController) Login2(ctx *gin.Context) {
	//  Lấy thông tin từ request
	var adminInput user_domain.SignIn

	if err := ctx.ShouldBindJSON(&adminInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	// Tìm kiếm admin trong cơ sở dữ liệu
	admin, err := l.AdminUseCase.Login(ctx, adminInput.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Tìm kiếm user trong database
	user, err := l.UserUseCase.Login(ctx, adminInput.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

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

	// Generate token
	accessToken, err = internal.CreateToken(l.Database.AccessTokenExpiresIn, user.ID, l.Database.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error()},
		)
		return
	}

	refreshToken, err = internal.CreateToken(l.Database.RefreshTokenExpiresIn, user.ID, l.Database.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error()},
		)
		return
	}

	ctx.SetCookie("access_token", accessToken, l.Database.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, l.Database.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", l.Database.AccessTokenMaxAge*60, "/", "localhost", false, false)

	// Trả về thông báo login thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "login successful",
	})
}
