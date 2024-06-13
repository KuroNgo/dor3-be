package vocabulary_controller

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (v *VocabularyController) UpdateOneVocabulary(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := v.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var vocabularyInput vocabulary_domain.Input
	if err := ctx.ShouldBindJSON(&vocabularyInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	updateVocabulary := vocabulary_domain.Vocabulary{
		Word:          vocabularyInput.Word,
		PartOfSpeech:  vocabularyInput.PartOfSpeech,
		Pronunciation: vocabularyInput.Pronunciation,
		ExplainEng:    vocabularyInput.ExplainEng,
		ExplainVie:    vocabularyInput.ExplainVie,
		ExampleVie:    vocabularyInput.ExampleVie,
		ExampleEng:    vocabularyInput.ExampleEng,
		FieldOfIT:     vocabularyInput.FieldOfIT,
		LinkURL:       vocabularyInput.LinkURL,
		UnitID:        vocabularyInput.UnitID,
		WhoUpdates:    admin.FullName,
	}

	data, err := v.VocabularyUseCase.UpdateOneInAdmin(ctx, &updateVocabulary)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
