package exam_answer

import (
	exam_answer_domain "clean-architecture/domain/exam_answer"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (e *ExamAnswerController) CreateOneExam(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	userID, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%s", currentUser))

	questionID := ctx.Query("question_id")
	idQuestion, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%s", questionID))

	var answerInput exam_answer_domain.Input
	if err := ctx.ShouldBindJSON(&answerInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	exam := exam_answer_domain.ExamAnswer{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		QuestionID:  idQuestion,
		Content:     answerInput.Content,
		SubmittedAt: time.Now(),
	}

	err := e.ExamAnswerUseCase.CreateOne(ctx, &exam)
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
