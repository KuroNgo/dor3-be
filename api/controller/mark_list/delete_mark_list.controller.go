package marklist_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

// DeleteOneMarkListByUser chỉ user mới có thể xóa
func (m *MarkListController) DeleteOneMarkListByUser(ctx *gin.Context) {
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

	markListID := ctx.Query("_id")

	err = m.MarkListUseCase.DeleteOneByUser(ctx, user.ID, markListID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	data, err := m.MarkVocabularyUseCase.FetchManyByMarkListID(ctx, markListID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, elMarkVocab := range data {
			err = m.MarkVocabularyUseCase.DeleteOne(ctx, elMarkVocab.ID.Hex())
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}
		}
	}()

	wg.Wait()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
