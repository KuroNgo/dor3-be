package exercise_options_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	Content    string `bson:"content" json:"content"`
	BlankIndex int    `bson:"blank_index" json:"blank_index"` // Chỉ số của từ cần điền vào câu, nếu là loại bài tập điền từ
}

type IExerciseOptionUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	UpdateOne(ctx context.Context, exerciseOptions *ExerciseOptions) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, exerciseOptions *ExerciseOptions) error
	DeleteOne(ctx context.Context, optionsID string) error
}
