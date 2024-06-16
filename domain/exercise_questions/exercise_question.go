package exercise_questions_domain

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionExerciseQuestion = "exercise_question"
)

type ExerciseQuestion struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	ExerciseID   primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`

	Content       string   `bson:"content" json:"content"`
	Type          string   `bson:"type" json:"type"`
	Level         int      `bson:"level" json:"level"`
	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type ExerciseQuestionResponse struct {
	ID         primitive.ObjectID           `bson:"_id" json:"_id"`
	ExerciseID primitive.ObjectID           `bson:"exercise_id" json:"exercise_id"`
	Vocabulary vocabulary_domain.Vocabulary `bson:"vocabulary" json:"vocabulary"`

	Content       string   `bson:"content" json:"content"`
	Type          string   `bson:"type" json:"type"`
	Level         int      `bson:"level" json:"level"`
	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type Response struct {
	Page                     int64                      `bson:"page" json:"page"`
	CurrentPage              int                        `bson:"current_page" json:"current_page"`
	Statistics               Statistics                 `bson:"statistics" json:"statistics"`
	ExerciseQuestionResponse []ExerciseQuestionResponse `json:"exercise_question" bson:"exercise_question"`
}

type Statistics struct {
	Count int64 `bson:"count" json:"count"`
}

type IExerciseQuestionRepository interface {
	FetchManyInAdmin(ctx context.Context, page string) (Response, error)
	FetchByIDInAdmin(ctx context.Context, id string) (ExerciseQuestionResponse, error)
	FetchManyByExerciseIDInAdmin(ctx context.Context, exerciseID string) (Response, error)
	FetchOneByExerciseIDInAdmin(ctx context.Context, exerciseID string) (ExerciseQuestionResponse, error)

	UpdateOneInAdmin(ctx context.Context, exerciseQuestion *ExerciseQuestion) (*mongo.UpdateResult, error)
	CreateOneInAdmin(ctx context.Context, exerciseQuestion *ExerciseQuestion) error
	DeleteOneInAdmin(ctx context.Context, exerciseID string) error
}
