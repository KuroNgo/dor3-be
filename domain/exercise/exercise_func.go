package exercise_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Input struct {
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`
	Title        string             `bson:"title" json:"title"`
	Content      string             `bson:"content" json:"content"`
	Type         string             `bson:"type" json:"type"` // Loại bài tập: ví dụ: trắc nghiệm, điền từ, sắp xếp câu, v.v.
	//Options      []string           `bson:"options" json:"options"`         // Các lựa chọn cho bài tập trắc nghiệm, nếu có
	CorrectAns string `bson:"correct_ans" json:"correct_ans"` // Đáp án đúng cho bài tập, nếu có
	BlankIndex int    `bson:"blank_index" json:"blank_index"` // Chỉ số của từ cần điền vào câu, nếu là loại bài tập điền từ
}

type IExerciseUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	UpdateOne(ctx context.Context, exerciseID string, exercise Exercise) error
	CreateOne(ctx context.Context, exercise *Exercise) error
	UpsertOne(ctx context.Context, id string, exercise *Exercise) (Response, error)
	DeleteOne(ctx context.Context, exerciseID string) error
}
