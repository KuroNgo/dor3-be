package vocabulary_controller

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/internal/cloud/cloudinary"
	"clean-architecture/internal/cloud/google"
	file_internal "clean-architecture/internal/file"
	"clean-architecture/internal/file/excel"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

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

	unitId := ctx.Request.FormValue("unit_id")
	idUnit, err := primitive.ObjectIDFromHex(unitId)

	word := ctx.Request.FormValue("word")
	partOfSpeech := ctx.Request.FormValue("part_of_speech")
	pronunciation := ctx.Request.FormValue("pronunciation")
	mean := ctx.Request.FormValue("mean")
	exampleVie := ctx.Request.FormValue("example_vie")
	exampleEng := ctx.Request.FormValue("example_eng")
	explainVie := ctx.Request.FormValue("explain_vie")
	explainEng := ctx.Request.FormValue("explain_eng")
	fieldOfIt := ctx.Request.FormValue("field_of_it")
	linkUrl := ctx.Request.FormValue("link_url")

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	if !file_internal.IsImage(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file format. Only images are allowed.",
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error opening uploaded file",
		})
		return
	}
	defer f.Close()

	// Tải file lên Cloudinary
	result, err := cloudinary.UploadImageToCloudinary(f, file.Filename, v.Database.CloudinaryUploadFolderUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	vocabularyRes := &vocabulary_domain.Vocabulary{
		Id:            primitive.NewObjectID(),
		Word:          word,
		PartOfSpeech:  partOfSpeech,
		Pronunciation: pronunciation,
		Mean:          mean,
		ExplainVie:    explainVie,
		ExampleEng:    exampleEng,
		ExplainEng:    explainEng,
		ExampleVie:    exampleVie,
		FieldOfIT:     fieldOfIt,
		LinkURL:       linkUrl,
		ImageURL:      result.ImageURL,
		AssetURL:      result.AssetID,
		UnitID:        idUnit,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		WhoUpdates:    admin.FullName,
	}

	err = v.VocabularyUseCase.CreateOneInAdmin(ctx, vocabularyRes)
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

	v.CreateAudioLatest(ctx)
	v.UploadAudioToCloudinary(ctx)
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
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, vocabulary := range result {
			unitID, err := v.UnitUseCase.FindUnitIDByUnitLevelInAdmin(ctx, vocabulary.UnitLevel, vocabulary.FieldOfIT)
			if err != nil {
				ctx.JSON(500, gin.H{"error": err.Error()})
				return
			}

			vocabularyTrimSpace := strings.ReplaceAll(vocabulary.Word, " ", "")

			vocab := vocabulary_domain.Vocabulary{
				Id:            primitive.NewObjectID(),
				UnitID:        unitID,
				Word:          vocabulary.Word,
				WordForConfig: vocabularyTrimSpace,
				PartOfSpeech:  vocabulary.PartOfSpeech,
				Pronunciation: vocabulary.Pronunciation,
				Mean:          vocabulary.Example,
				ExampleEng:    vocabulary.ExampleEng,
				ExampleVie:    vocabulary.ExampleVie,
				ExplainEng:    vocabulary.ExplainEng,
				ExplainVie:    vocabulary.ExplainVie,
				FieldOfIT:     vocabulary.FieldOfIT,
				LinkURL:       "",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				WhoUpdates:    admin.FullName,
			}
			vocabularies = append(vocabularies, vocab)
		}
	}()

	wg.Wait()

	successCount := 0
	for _, vocabulary := range vocabularies {
		err = v.VocabularyUseCase.CreateOneByNameUnitInAdmin(ctx, &vocabulary)
		if err != nil {
			continue
		}
		successCount++
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success process in create vocabulary with file",
		"success_count": successCount,
	})

	v.CreateAudioLatest(ctx)
	v.UploadAudioToCloudinary(ctx)
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
		unitID, err := v.UnitUseCase.FindUnitIDByUnitLevelInAdmin(ctx, vocabulary.UnitLevel, vocabulary.FieldOfIT)
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
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			WhoUpdates:    admin.FullName,
		}
		vocabularies = append(vocabularies, v)
	}

	successCount := 0
	for _, vocabulary := range vocabularies {
		err = v.VocabularyUseCase.CreateOneByNameUnitInAdmin(ctx, &vocabulary)
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

	data, err := v.VocabularyUseCase.GetAllVocabularyInAdmin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, audio := range data {
			_ = google.CreateTextToSpeech(audio)
		}
	}()

	wg.Wait()

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

	data, err := v.VocabularyUseCase.GetLatestVocabularyInAdmin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, audio := range data {
			_ = google.CreateTextToSpeech(audio)
		}
	}()
	wg.Wait()

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, audioFileName := range files {
			audio := strings.TrimSuffix(audioFileName, ".mp3")

			// Mở từng tệp
			f, err := os.Open(filepath.Join(dir, audioFileName))
			if err != nil {
				// Đảm bảo đóng tệp nếu có lỗi
				f.Close()
				return
			}

			// Kiểm tra xem file có phải là MP3 không
			if !file_internal.IsMP3(audioFileName) {
				f.Close() // Đóng tệp trước khi trả về lỗi
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("%s is not an MP3 file", audioFileName),
				})
				return
			}

			// Đóng tệp ngay sau khi sử dụng
			defer f.Close()

			// Upload file lên Cloudinary
			data, err := cloudinary.UploadAudioToCloudinary(f, audioFileName, v.Database.CloudinaryUploadFolderAudioVocabulary)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Error uploading file %s to Cloudinary: %s", audioFileName, err),
				})
				return
			}

			// Tìm ID của từ vựng dựa trên tên file
			wordID, err := v.VocabularyUseCase.FindVocabularyIDByVocabularyConfigInAdmin(ctx, audio)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Error finding vocabulary ID for file %s: %s", audioFileName, err.Error()),
				})
				return
			}

			vocabularyRes := &vocabulary_domain.Vocabulary{
				Id:      wordID,
				LinkURL: data.AudioURL,
			}

			err = v.VocabularyUseCase.UpdateOneAudioInAdmin(ctx, vocabularyRes)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Error updating vocabulary for file %s: %s", audioFileName, err.Error()),
				})
				return
			}
		}
	}()

	wg.Wait()

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	v.DeleteFolderOfVocabulary(ctx)
}
