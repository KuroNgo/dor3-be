package vocabulary_controller

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (v *VocabularyController) FetchByWord(ctx *gin.Context) {
	var word vocabulary_domain.FetchByWordInput

	if err := ctx.ShouldBindJSON(&word); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	vocabulary, err := v.VocabularyUseCase.FetchByWord(ctx, word.Word)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"word": vocabulary,
		},
	})
}

func (v *VocabularyController) FetchByLesson(ctx *gin.Context) {
	var lesson vocabulary_domain.FetchByLessonInput
	if err := ctx.ShouldBindJSON(&lesson); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	vocabulary, err := v.VocabularyUseCase.FetchByWord(ctx, lesson.Lesson)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"word": vocabulary,
		},
	})
}
