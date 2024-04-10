package mark_vocabulary_controller

import (
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (m *MarkVocabularyController) CreateOneMarkVocabulary(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	user, err := m.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	markListID := ctx.Query("mark_list_id")
	markListData, err := primitive.ObjectIDFromHex(markListID)

	vocabularyID := ctx.Query("vocabulary_id")
	vocabularyData, err := primitive.ObjectIDFromHex(vocabularyID)

	markVocabularyReq := mark_vocabulary_domain.MarkToFavourite{
		ID:           primitive.NewObjectID(),
		UserId:       user.ID,
		MarkListID:   markListData,
		VocabularyID: vocabularyData,
	}

	err = m.MarkVocabularyUseCase.CreateOne(ctx, &markVocabularyReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error create mark vocabulary",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
