package quiz_question_domain

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionQuizQuestion = "quiz_question"
)

type QuizQuestion struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	QuizID       primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
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

type QuizQuestionResponse struct {
	ID         primitive.ObjectID           `bson:"_id" json:"_id"`
	QuizID     primitive.ObjectID           `bson:"quiz_id" json:"quiz_id"`
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
	Page         int64                  `bson:"page" json:"page"`
	CurrentPage  int64                  `bson:"current_page" json:"current_page"`
	Statistics   Statistics             `bson:"statistics" json:"statistics"`
	QuizQuestion []QuizQuestionResponse `json:"quiz_question" bson:"quiz_question"`
}

type Statistics struct {
	Count int64 `bson:"count" json:"count"`
}

type IQuizQuestionRepository interface {
	FetchManyInAdmin(ctx context.Context, page string) (Response, error)
	FetchByIDInAdmin(ctx context.Context, id string) (QuizQuestion, error)
	FetchManyByQuizIDInAdmin(ctx context.Context, quizID string) (Response, error)
	FetchOneByQuizIDInAdmin(ctx context.Context, quizID string) (QuizQuestion, error)

	UpdateOneInAdmin(ctx context.Context, quizQuestion *QuizQuestion) (*mongo.UpdateResult, error)
	CreateOneInAdmin(ctx context.Context, quizQuestion *QuizQuestion) error
	DeleteOneInAdmin(ctx context.Context, quizID string) error
}
