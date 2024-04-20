package mark_vocabulary_controller

import (
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (m *MarkVocabularyController) CreateOneMarkVocabulary(ctx *gin.Context) {
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
			"message": user.FullName + "You are not authorized to perform this action!",
		})
		return
	}

	var input mark_vocabulary_domain.Input
	if err = ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	markListID, err := primitive.ObjectIDFromHex(input.MarkListID)
	vocabularyID, err := primitive.ObjectIDFromHex(input.VocabularyID)

	markVocabularyReq := mark_vocabulary_domain.MarkToFavourite{
		ID:           primitive.NewObjectID(),
		UserId:       user.ID,
		MarkListID:   markListID,
		VocabularyID: vocabularyID,
	}

	err = m.MarkVocabularyUseCase.CreateOne(ctx, &markVocabularyReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error create mark vocabulary",
			"message": err.Error(),
		})
		return
	}

	err = m.VocabularyUseCase.UpdateIsFavourite(ctx, input.VocabularyID, 1)
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
