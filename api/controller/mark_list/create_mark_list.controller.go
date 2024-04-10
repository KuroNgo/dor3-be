package marklist_controller

import (
	mark_list_domain "clean-architecture/domain/mark_list"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (m *MarkListController) CreateOneMarkList(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	user, err := m.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
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
		ID:          primitive.NewObjectID(),
		UserID:      user.ID,
		NameList:    markListInput.NameList,
		Description: markListInput.Description,
		WhoCreated:  user.FullName,
	}

	err = m.MarkListUseCase.CreateOne(ctx, &markListReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error create list mark vocabulary",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
