package marklist_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOneMarkList chỉ user mới có thể xóa
func (m *MarkListController) DeleteOneMarkList(ctx *gin.Context) {
	markListID := ctx.Query("_id")

	err := m.MarkListUseCase.DeleteOne(ctx, markListID)
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
