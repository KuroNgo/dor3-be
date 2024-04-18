package vocabulary_controller

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/internal"
	"clean-architecture/internal/cloud/cloudinary"
	"clean-architecture/internal/cloud/google"
	file_internal "clean-architecture/internal/file"
	"clean-architecture/internal/file/excel"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

func (v *VocabularyController) CreateOneVocabulary(ctx *gin.Context) {
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
		Mean:          vocabularyInput.Mean,
		ExplainVie:    vocabularyInput.ExplainVie,
		ExampleEng:    vocabularyInput.ExampleEng,
		ExplainEng:    vocabularyInput.ExplainEng,
		ExampleVie:    vocabularyInput.ExampleVie,
		FieldOfIT:     vocabularyInput.FieldOfIT,
		LinkURL:       vocabularyInput.LinkURL,
		UnitID:        vocabularyInput.UnitID,
		IsFavourite:   0,
		WhoUpdates:    admin.FullName,
	}

	err = v.VocabularyUseCase.CreateOne(ctx, vocabularyRes)
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

func (v *VocabularyController) CreateVocabularyWithFileInAdmin(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := v.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	err = ctx.Request.ParseMultipartForm(8 << 20) // 8MB max size
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
		unitID, err := v.VocabularyUseCase.FindUnitIDByUnitLevel(ctx, vocabulary.UnitLevel)
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
			Mean:          vocabulary.Example,
			ExampleEng:    vocabulary.ExampleEng,
			ExampleVie:    vocabulary.ExampleVie,
			ExplainEng:    vocabulary.ExplainEng,
			ExplainVie:    vocabulary.ExplainVie,
			FieldOfIT:     vocabulary.FieldOfIT,
			LinkURL:       "",
			IsFavourite:   0,
			WhoUpdates:    admin.FullName,
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

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"success_count": successCount,
	})
}

func (v *VocabularyController) CreateVocabularyWithFileInUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := v.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	err = ctx.Request.ParseMultipartForm(8 << 20) // 8MB max size
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
		unitID, err := v.VocabularyUseCase.FindUnitIDByUnitLevel(ctx, vocabulary.UnitLevel)
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
			Mean:          vocabulary.Example,
			ExampleEng:    vocabulary.ExampleEng,
			ExampleVie:    vocabulary.ExampleVie,
			ExplainEng:    vocabulary.ExplainEng,
			ExplainVie:    vocabulary.ExplainVie,
			FieldOfIT:     vocabulary.FieldOfIT,
			LinkURL:       "",
			IsFavourite:   0,
			WhoUpdates:    admin.FullName,
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

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"success_count": successCount,
	})
}

func (v *VocabularyController) CreateAudio(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := v.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	data, err := v.VocabularyUseCase.GetAllVocabulary(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, audio := range data {
		_ = google.CreateTextToSpeech(audio)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (v *VocabularyController) CreateAudioLatest(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := v.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	data, err := v.VocabularyUseCase.GetLatestVocabulary(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, audio := range data {
		_ = google.CreateTextToSpeech(audio)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (v *VocabularyController) UploadAudioToCloudinary(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := v.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	// Lấy danh sách tệp trong thư mục audio
	dir := "audio"
	files, err := google.ListFilesInDirectory(dir)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

		_, err = cloudinary.UploadToCloudinary(f, filename.(string), v.Database.CloudinaryUploadFolderAudioVocabulary)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		google.DeleteFile(audioFile.Filename)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

}
