package vocabulary_controller

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (v *VocabularyController) UpsertOneVocabulary(ctx *gin.Context) {
	vocabularyID := ctx.Query("_id")

	var vocabularyInput vocabulary_domain.Input
	if err := ctx.ShouldBindJSON(&vocabularyInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	upsertVocabulary := vocabulary_domain.Vocabulary{
		Word:          vocabularyInput.Word,
		PartOfSpeech:  vocabularyInput.PartOfSpeech,
		Pronunciation: vocabularyInput.Pronunciation,
		Example:       vocabularyInput.Example,
		FieldOfIT:     vocabularyInput.FieldOfIT,
		LinkURL:       vocabularyInput.LinkURL,
		UnitID:        vocabularyInput.UnitID,
	}

	vocabularyRes, err := v.VocabularyUseCase.UpsertOne(ctx, vocabularyID, &upsertVocabulary)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   vocabularyRes,
	})
}
