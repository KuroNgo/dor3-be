package exercise_answer_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Input struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Content     string    `bson:"content" json:"content"`
	BlankIndex  int       `bson:"blank_index" json:"blank_index"` // Chỉ số của từ cần điền vào câu, nếu là loại bài tập điền từ
	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type IExerciseAnswerUseCase interface {
	FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (Response, error)
	CreateOne(ctx context.Context, exerciseAnswer *ExerciseAnswer) error
	DeleteOne(ctx context.Context, exerciseID string) error
}
