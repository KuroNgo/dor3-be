package vocabulary_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (v *VocabularyController) FetchByWord(ctx *gin.Context) {
	word := ctx.Query("word")

	vocabulary, err := v.VocabularyUseCase.FetchByWordInBoth(ctx, word)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   vocabulary,
	})
}

func (v *VocabularyController) FetchByLesson(ctx *gin.Context) {
	lesson := ctx.Query("field_of_it")

	vocabulary, err := v.VocabularyUseCase.FetchByLessonInBoth(ctx, lesson)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   vocabulary,
	})
}

func (v *VocabularyController) FetchMany(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	// Truyền giá trị page từ người dùng vào use case
	vocabulary, err := v.VocabularyUseCase.FetchManyInBoth(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   vocabulary,
	})
}

func (v *VocabularyController) FetchByIdUnit(ctx *gin.Context) {
	idUnit := ctx.Query("unit_id")

	vocabulary, err := v.VocabularyUseCase.FetchByIdUnitInAdmin(ctx, idUnit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"vocabulary": vocabulary,
	})
}

func (v *VocabularyController) FetchByIdVocabulary(ctx *gin.Context) {
	idVocabulary := ctx.Query("_id")

	vocabulary, err := v.VocabularyUseCase.GetVocabularyByIdInAdmin(ctx, idVocabulary)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"vocabulary": vocabulary,
	})
}

func (v *VocabularyController) FetchAllVocabulary(ctx *gin.Context) {
	vocabulary, err := v.VocabularyUseCase.GetAllVocabularyInAdmin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"vocabulary": vocabulary,
	})
}
func (v *VocabularyController) FetchAllVocabularyLatest(ctx *gin.Context) {
	vocabulary, err := v.VocabularyUseCase.GetLatestVocabularyInAdmin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"vocabulary": vocabulary,
	})
}
