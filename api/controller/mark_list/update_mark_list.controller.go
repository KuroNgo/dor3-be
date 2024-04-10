package marklist_controller

import (
	mark_list_domain "clean-architecture/domain/mark_list"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateOneMarkList update one mark list for user
func (m *MarkListController) UpdateOneMarkList(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	user, err := m.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	markListID := ctx.Query("_id")
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

	err = m.MarkListUseCase.UpdateOne(ctx, markListID, markListReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}