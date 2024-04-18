package exercise_options_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionExerciseOptions = "exercise_options"
)

type ExerciseOptions struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Content    string `bson:"content" json:"content"`
	BlankIndex int    `bson:"blank_index" json:"blank_index"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt  time.Time `bson:"update_at" json:"update_at"`
	WhoUpdate string    `bson:"who_update" json:"who_update"`
}

type Response struct {
	ExerciseOptions []ExerciseOptions
}

type IExerciseOptionRepository interface {
	FetchManyByQuestionID(ctx context.Context, questionID string) (Response, error)
	UpdateOne(ctx context.Context, exerciseOptions *ExerciseOptions) (*mongo.UpdateResult, error)
	CreateOne(ctx context.Context, exerciseOptions *ExerciseOptions) error
	DeleteOne(ctx context.Context, optionsID string) error
}
