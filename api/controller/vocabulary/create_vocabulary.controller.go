package vocabulary_controller

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (v *VocabularyController) CreateOneLesson(ctx *gin.Context) {
	var vocabularyInput vocabulary_domain.Input

	if err := ctx.ShouldBindJSON(&vocabularyInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := internal.IsValidVocabulary(vocabularyInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	vocabularyRes := &vocabulary_domain.Vocabulary{
		Id:            primitive.NewObjectID(),
		Word:          vocabularyInput.Word,
		PartOfSpeech:  vocabularyInput.PartOfSpeech,
		Pronunciation: vocabularyInput.Pronunciation,
		Example:       vocabularyInput.Example,
		FieldOfIT:     vocabularyInput.FieldOfIT,
		LinkURL:       vocabularyInput.LinkURL,
		LessonID:      vocabularyInput.LessonID,
	}

	err := v.VocabularyUseCase.CreateOne(ctx, vocabularyRes)
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
