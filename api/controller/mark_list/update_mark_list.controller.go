package marklist_controller

import (
	mark_list_domain "clean-architecture/domain/mark_list"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateOneMarkList update one mark list for user
func (m *MarkListController) UpdateOneMarkList(ctx *gin.Context) {
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

	var markListInput mark_list_domain.Input
	if err = ctx.ShouldBindJSON(&markListInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	markListReq := mark_list_domain.MarkList{
		UserID:      user.ID,
		NameList:    markListInput.NameList,
		Description: markListInput.Description,
		WhoCreated:  user.FullName,
	}

	data, err := m.MarkListUseCase.UpdateOne(ctx, &markListReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
