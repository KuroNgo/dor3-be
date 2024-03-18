package admin_controller

import (
	admin_domain "clean-architecture/domain/admin"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func (a *AdminController) SignUp(ctx *gin.Context) {
	var admin admin_domain.SignUp
	if err := ctx.ShouldBindJSON(&admin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	if !internal.EmailValid(admin.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email không hợp lệ !",
		})
		return
	}
	// Bên phía client sẽ phải so sánh password thêm một lần nữa đã đúng chưa
	if !internal.PasswordStrong(admin.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Mật khẩu phải có ít nhất 8 ký tự, bao gồm chữ hoa, chữ thường và số !",
		})
		return
	}

	// Băm mật khẩu
	hashedPassword, err := internal.HashPassword(admin.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	admin.Password = hashedPassword
	admin.Password = internal.Santize(admin.Password)
	admin.Email = internal.Santize(admin.Email)
	admin.Email = strings.ToLower(admin.Email)

	newAdmin := admin_domain.Admin{
		Address:   admin.Address,
		FullName:  admin.FullName,
		Avatar:    admin.Avatar,
		Email:     admin.Email,
		Password:  hashedPassword,
		Role:      "admin",
		Phone:     admin.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// thực hiện đăng ký người dùng
	err = a.AdminUseCase.CreateOne(ctx, newAdmin)
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
