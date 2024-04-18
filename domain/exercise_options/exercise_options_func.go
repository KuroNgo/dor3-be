package exercise_options_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`
	Content    string             `bson:"content" json:"content"`
	BlankIndex int                `bson:"blank_index" json:"blank_index"` // Chỉ số của từ cần điền vào câu, nếu là loại bài tập điền từ
}

type IExerciseOptionUseCase interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	UpdateOne(ctx context.Context, exerciseOptions *ExerciseOptions) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, exerciseOptions *ExerciseOptions) error
	DeleteOne(ctx context.Context, optionsID string) error
}
