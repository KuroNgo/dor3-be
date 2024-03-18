package user_controller

import (
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (u *UserController) SignUp(ctx *gin.Context) {
	var user user_domain.SignUp
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	if !internal.EmailValid(user.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email không hợp lệ !",
		})
		return
	}
	// Bên phía client sẽ phải so sánh password thêm một lần nữa đã đúng chưa
	if !internal.PasswordStrong(user.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Mật khẩu phải có ít nhất 8 ký tự, bao gồm chữ hoa, chữ thường và số !",
		})
		return
	}

	// Băm mật khẩu
	hashedPassword, err := internal.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	user.Password = hashedPassword
	user.Password = internal.Santize(user.Password)
	user.Email = internal.Santize(user.Email)

	newUser := user_domain.User{
		ID:         primitive.NewObjectID(),
		FullName:   user.FullName,
		AvatarURL:  user.AvatarURL,
		Specialize: user.Specialize,
		Email:      user.Email,
		Password:   hashedPassword,
		Verified:   false,
		Provider:   "fe-it",
		Role:       "user",
		Phone:      user.Phone,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// thực hiện đăng ký người dùng
	err = u.UserUseCase.Create(ctx, newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
