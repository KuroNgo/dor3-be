package marklist_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *MarkListController) FetchManyMarkListByUserID(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	user, err := m.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))

	markList, err := m.MarkListUseCase.FetchManyByUserID(ctx, fmt.Sprint(user.ID))
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
