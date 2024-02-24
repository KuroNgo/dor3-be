package audio_controller

import (
	audio_domain "clean-architecture/domain/audio"
	quiz_domain "clean-architecture/domain/quiz"
	file_internal "clean-architecture/internal/file"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"path/filepath"
)

// CreateAudioInFireBaseAndSaveMetaDataInDatabase used for
// collect information from file and upload to firebase with audio and database with metadata
func (au *AudioController) CreateAudioInFireBaseAndSaveMetaDataInDatabase(ctx *gin.Context) {
	var quiz quiz_domain.Quiz
	if err := ctx.ShouldBindJSON(&quiz); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Just filename contains "file" in OS (operating system)
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Define the path where the file will be saved
	filePath := filepath.Join("uploads", file.Filename)

	// get filename mp3
	filename, err := file_internal.GetNameFileMP3(filePath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get time duration mp3
	duration, err := file_internal.GetDurationFileMP3(filePath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ID := ctx.PostForm("quiz_id")
	quizID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		fmt.Println("Fail to convert: ", err)
		return
	}

	// the metadata will be saved in database
	metadata := &audio_domain.AutoMatch{
		QuizID:        quizID,
		Filename:      filename,
		AudioDuration: duration,
	}

	// save data in database
	err = au.AudioUseCase.CreateOne(ctx, metadata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// process in firebase below code

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
