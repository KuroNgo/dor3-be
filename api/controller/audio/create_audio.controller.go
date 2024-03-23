package audio_controller

import (
	audio_domain "clean-architecture/domain/audio"
	file_internal "clean-architecture/internal/file"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

// CreateAudioInFireBaseAndSaveMetaDataInDatabase used for
// collect information from file and upload to firebase with audio and database with metadata
func (au *AudioController) CreateAudioInFireBaseAndSaveMetaDataInDatabase(ctx *gin.Context) {
	// Just filename contains "file" in OS (operating system)
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !file_internal.IsMP3(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// the metadata will be saved in database
	metadata := &audio_domain.Audio{
		Id: primitive.NewObjectID(),
		//QuizID:        quizID,
		Filename: file.Filename,
		Size:     file.Size,
	}

	// save data in database
	err = au.AudioUseCase.CreateOne(ctx, metadata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// process in firebase below code

	// Xóa file sau khi đã sử dụng
	err = os.Remove("./" + file.Filename)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
