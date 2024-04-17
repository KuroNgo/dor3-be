package marklist_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOneMarkList chỉ user mới có thể xóa
func (m *MarkListController) DeleteOneMarkList(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	user, err := m.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": user.FullName + "You are not authorized to perform this action!",
		})
		return
	}

	markListID := ctx.Query("_id")

	err = m.MarkListUseCase.DeleteOne(ctx, markListID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
