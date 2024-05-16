package admin_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *AdminController) FetchManyUser(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := a.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	data, err := a.UserUseCase.FetchMany(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   data,
	})
}

func (a *AdminController) FetchUserByID(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := a.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	userID := ctx.Query("user_id")
	data, err := a.UserUseCase.GetByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   data,
	})
}
