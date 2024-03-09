package quiz_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionQuiz = "quiz"
)

type Quiz struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Question      string             `bson:"question" json:"question"`
	Options       []string           `bson:"options" json:"options"`
	CorrectAnswer string             `bson:"correct_answer" json:"correct_answer"`
	Explanation   string             `bson:"explanation" json:"explanation"`
	QuestionType  string             `bson:"question_type" json:"question_type"`

	// admin add metadata of file and system will be found it
	Filename      string `bson:"filename" json:"filename"`
	AudioDuration string `bson:"audio_duration" json:"audio_duration"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	Question      string   `bson:"question" json:"question"`
	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`
	Explanation   string   `bson:"explanation" json:"explanation"`

	// QuestionType can be included checkbox, check radius or write correct answer
	QuestionType string `bson:"question_type" json:"question_type"`

	// admin add metadata of file and system will be found it
	Filename      string `bson:"filename" json:"filename"`
	AudioDuration string `bson:"audio_duration" json:"audio_duration"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

//go:generate mockery --name IQuizRepository
type IQuizRepository interface {
	FetchByID(ctx context.Context, quizID string) (*Quiz, error)
	FetchMany(ctx context.Context) ([]Quiz, error)
	FetchToDeleteMany(ctx context.Context) (*[]Quiz, error)
	UpdateOne(ctx context.Context, quizID string, quiz Quiz) error
	CreateOne(ctx context.Context, quiz *Quiz) error
	UpsertOne(c context.Context, id string, quiz *Quiz) (*Response, error)
	DeleteOne(ctx context.Context, quizID string) error
}
