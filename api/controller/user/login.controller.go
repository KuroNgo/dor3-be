package controller

import (
	"clean-architecture/bootstrap"
	user_domain "clean-architecture/domain/request/user"
	"clean-architecture/domain/response"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginController struct {
	LoginUseCase user_domain.ILoginUseCase
	Env          *bootstrap.Database
}

func (lc *LoginController) LoginByUserName(ctx *gin.Context) {
	var request user_domain.LoginUsernameRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	user, err := lc.LoginUseCase.GetUserByUsername(ctx, request.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: "User not found with given username",
		})
		return
	}

	if err = internal.VerifyPassword(user.Password, request.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "Invalid credentials"})
		return
	}

	// Thiếu Token (Missing Token)
}

func (lc *LoginController) LoginByEmail(ctx *gin.Context) {
	var request user_domain.LoginEmailRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	user, err := lc.LoginUseCase.GetUserByEmail(ctx, request.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: "User not found with given email",
		})
		return
	}

	if err = internal.VerifyPassword(user.Password, request.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Message: "Invalid Credentials",
		})
		return
	}

	// Missing Token
}
