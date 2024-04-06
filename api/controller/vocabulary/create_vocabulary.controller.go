package vocabulary_controller

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/internal"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"clean-architecture/internal/file/excel"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

func (v *VocabularyController) CreateOneVocabulary(ctx *gin.Context) {
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
		UnitID:        vocabularyInput.UnitID,
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

func (v *VocabularyController) CreateVocabularyWithFile(ctx *gin.Context) {
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

	result, err := excel.ReadFileForVocabulary(file.Filename)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var vocabularies []vocabulary_domain.Vocabulary

	for _, vocabulary := range result {
		unitID, err := v.VocabularyUseCase.FindUnitIDByUnitName(ctx, vocabulary.UnitID)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		v := vocabulary_domain.Vocabulary{
			Id:            primitive.NewObjectID(),
			UnitID:        unitID,
			Word:          vocabulary.Word,
			PartOfSpeech:  vocabulary.PartOfSpeech,
			Pronunciation: vocabulary.Pronunciation,
			Example:       vocabulary.Example,
			FieldOfIT:     vocabulary.FieldOfIT,
			LinkURL:       "",
		}
		vocabularies = append(vocabularies, v)
	}

	successCount := 0
	for _, vocabulary := range vocabularies {
		err = v.VocabularyUseCase.CreateOneByNameUnit(ctx, &vocabulary)
		if err != nil {
			continue
		}
		successCount++
	}

	if successCount == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create any vocabulary",
			"message": "Any value have exist in database",
		})
		return
	}

	// Tạo audio từ từng
	vocabulary, err := v.VocabularyUseCase.GetAllVocabulary(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, vocabularyElement := range vocabulary {
		err = file_internal.CreateTextToSpeech(vocabularyElement)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Lấy danh sách tệp trong thư mục audio
	dir := "audio"
	files, err := file_internal.ListFilesInDirectory(dir)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var vocabularies2 []vocabulary_domain.Vocabulary
	for _, audioFile := range files {
		f, err := audioFile.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error opening uploaded file",
			})
			return
		}

		if !file_internal.IsMP3(audioFile.Filename) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Not an MP3 file",
			})
			return
		}

		filename, ok := ctx.Get("filePath")
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "filename not found",
			})
			return
		}

		res, err := cloudinary.UploadToCloudinary(f, filename.(string), v.Database.CloudinaryUploadFolderAudioVocabulary)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		vu := vocabulary_domain.Vocabulary{
			LinkURL: res.AssetID,
		}

		vocabularies2 = append(vocabularies2, vu)
		file_internal.DeleteFile(audioFile.Filename)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"success_count": successCount,
	})
}
