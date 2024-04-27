package admin_controller

import (
	admin_domain "clean-architecture/domain/admin"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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
			"message": "Email Invalid !",
		})
		return
	}

	if !internal.PasswordStrong(admin.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Mật khẩu phải có ít nhất 8 ký tự, bao gồm chữ hoa, chữ thường và số !",
		})
		return
	}

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

	newAdmin := admin_domain.Admin{
		Id:        primitive.NewObjectID(),
		Address:   admin.Address,
		FullName:  admin.FullName,
		AvatarURL: admin.Avatar,
		Email:     admin.Email,
		Password:  hashedPassword,
		Role:      "admin",
		Phone:     admin.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = a.AdminUseCase.CreateOne(ctx, newAdmin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
