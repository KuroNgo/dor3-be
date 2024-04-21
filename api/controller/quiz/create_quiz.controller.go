package quiz_controller

import (
	quiz_domain "clean-architecture/domain/quiz"
	"clean-architecture/internal"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (q *QuizController) CreateOneQuiz(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	user, err := q.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

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

	quizRes := &quiz_domain.Quiz{
		ID:           primitive.NewObjectID(),
		LessonID:     quizInput.LessonID,
		UnitID:       quizInput.UnitID,
		VocabularyID: quizInput.VocabularyID,
		Title:        quizInput.Title,
		Description:  quizInput.Description,
		Duration:     quizInput.Duration,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		WhoUpdates:   user.FullName,
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

func (q *QuizController) CreateFromFileQuiz(ctx *gin.Context) {

}
