package exercise_answer_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Input struct {
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Content    string `bson:"content" json:"content"`
	BlankIndex int    `bson:"blank_index" json:"blank_index"` // Chỉ số của từ cần điền vào câu, nếu là loại bài tập điền từ

	SubmittedAt time.Time `bson:"submitted_at" json:"submitted_at"`
}

type IExerciseAnswerUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	CreateOne(ctx context.Context, exerciseAnswer *ExerciseAnswer) error
	DeleteOne(ctx context.Context, exerciseID string) error
}
