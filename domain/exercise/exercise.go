package exercise_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionExercise = "exercise"
)

type Exercise struct {
	Id         primitive.ObjectID `bson:"_id" json:"_id"`
	Vocabulary primitive.ObjectID `bson:"vocabulary" json:"vocabulary"`
	Title      string             `bson:"title" json:"title"`
	Content    string             `bson:"content" json:"content"`
	Type       string             `bson:"type" json:"type"`               // Loại bài tập: ví dụ: trắc nghiệm, điền từ, sắp xếp câu, v.v.
	Options    []string           `bson:"options" json:"options"`         // Các lựa chọn cho bài tập trắc nghiệm, nếu có
	CorrectAns string             `bson:"correct_ans" json:"correct_ans"` // Đáp án đúng cho bài tập, nếu có
	BlankIndex []int              `bson:"blank_index" json:"blank_index"` // Chỉ số của từ cần điền vào câu, nếu là loại bài tập điền từ
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	WhoUpdates string             `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	Exercise []Exercise `bson:"data" json:"data"`
	Count    int        `bson:"count" json:"count"`
}

type IExerciseRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	UpdateOne(ctx context.Context, exerciseID string, exercise Exercise) error
	CreateOne(ctx context.Context, exercise *Exercise) error
	UpsertOne(ctx context.Context, id string, exercise *Exercise) (Response, error)
	DeleteOne(ctx context.Context, exerciseID string) error
}
