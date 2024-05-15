package marklist_controller

import (
	"clean-architecture/internal"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteOneMarkList chỉ user mới có thể xóa
func (m *MarkListController) DeleteOneMarkList(ctx *gin.Context) {
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

	err = m.MarkListUseCase.DeleteOne(ctx, markListID)
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

	internal.Wg.Add(1)
	go func() {
		defer internal.Wg.Done()
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

	internal.Wg.Wait()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
