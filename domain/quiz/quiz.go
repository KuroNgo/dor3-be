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
	Level         int                `bson:"level" json:"level"`

	// admin add metadata of file and system will be found it
	Filename      string `bson:"filename" json:"filename"`
	AudioDuration string `bson:"audio_duration" json:"audio_duration"`

	IsComplete int       `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	Quiz []Quiz `bson:"data" json:"data"`
}

//go:generate mockery --name IQuizRepository
type IQuizRepository interface {
	//FetchTenQuizButEnoughAllSkill(ctx context.Context) ([]Response, error)
	FetchMany(ctx context.Context) (Response, error)
	UpdateOne(ctx context.Context, quizID string, quiz Quiz) error
	CreateOne(ctx context.Context, quiz *Quiz) error
	UpsertOne(c context.Context, id string, quiz *Quiz) (Response, error)
	DeleteOne(ctx context.Context, quizID string) error
}
