package quiz_controller

import (
	quiz_domain "clean-architecture/domain/quiz"
	"clean-architecture/internal"
	file_internal "clean-architecture/internal/file"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

func (q *QuizController) CreateOneQuiz(ctx *gin.Context) {
	var quizInput quiz_domain.Input

	if err := ctx.ShouldBindJSON(&quizInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := internal.IsValidQuiz(quizInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	// Just filename contains "file" in OS (operating system)
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Lưu file vào thư mục tạm thời
	err = ctx.SaveUploadedFile(file, "./"+file.Filename)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// get time duration mp3
	duration, err := file_internal.GetDurationFileMP3(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//ID := ctx.Query("quiz_id")
	//quizID, err := primitive.ObjectIDFromHex(ID)
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

	quizRes := &quiz_domain.Quiz{
		ID:            primitive.NewObjectID(),
		Question:      quizInput.Question,
		Options:       quizInput.Options,
		CorrectAnswer: quizInput.CorrectAnswer,
		Explanation:   quizInput.Explanation,
		QuestionType:  quizInput.QuestionType,
		Filename:      file.Filename,
		AudioDuration: duration,
	}

	err = q.QuizUseCase.CreateOne(ctx, quizRes)
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
