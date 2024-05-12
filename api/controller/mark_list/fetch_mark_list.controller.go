package marklist_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *MarkListController) FetchManyMarkListByUserID(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	user, err := m.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	markList, err := m.MarkListUseCase.FetchManyByUserID(ctx, user.ID.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"mark_list": markList,
	})
}

func (m *MarkListController) FetchManyMarkListByIDInAdmin(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := m.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	idMarkList := ctx.Query("_id")
	markList, err := m.MarkListUseCase.FetchById(ctx, idMarkList)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"mark_list": markList,
	})
}

func (m *MarkListController) FetchManyInAdmin(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := m.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	markList, err := m.MarkListUseCase.FetchMany(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"mark_list": markList,
	})
}
