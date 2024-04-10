package marklist_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Deprecated: FetchManyMarkList is deprecated, use FetchManyMarkListByUserId instead.
func (m *MarkListController) FetchManyMarkList(ctx *gin.Context) {
	markList, err := m.MarkListUseCase.FetchMany(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   markList,
	})
}
