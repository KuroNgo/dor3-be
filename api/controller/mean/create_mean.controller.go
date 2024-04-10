package mean_controller

import (
	mean_domain "clean-architecture/domain/mean"
	file_internal "clean-architecture/internal/file"
	"clean-architecture/internal/file/excel"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

// Deprecated: CreateMeanWithFile
func (m *MeanController) CreateMeanWithFile(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(8 << 20) // 8MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("files")
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !file_internal.IsExcel(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not an Excel file",
		})
		return
	}

	err = ctx.SaveUploadedFile(file, file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		err := os.Remove(file.Filename)
		if err != nil {
			fmt.Printf("Failed to delete temporary file: %v\n", err)
		}
	}()

	result, err := excel.ReadFileForMean(file.Filename)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var means []mean_domain.Mean
	for _, mean := range result {
		vocabularyID, err := m.MeanUseCase.FindVocabularyIDByWord(ctx, mean.VocabularyID)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		m := mean_domain.Mean{
			ID:           primitive.NewObjectID(),
			VocabularyID: vocabularyID,
			Description:  mean.ExplainEng,
			Example:      mean.ExampleEng,
			VietSub:      mean.ExplainVie,
			FieldOfIT:    mean.LessonID,
			SynonymID:    "",
			AntonymID:    "",
		}
		means = append(means, m)
	}

	successCount := 0
	for _, mean := range means {
		err = m.MeanUseCase.CreateOneByWord(ctx, &mean)
		if err != nil {
			continue
		}
		successCount++
	}

	if successCount == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create any lesson",
			"message": "Any value have exist in database",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"success_count": successCount,
	})
}
