package quiz_domain

import "context"

type Input struct {
	Question      string   `bson:"question" json:"question"`
	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`
	Explanation   string   `bson:"explanation" json:"explanation"`

	// QuestionType can be included checkbox, check radius or write correct answer
	QuestionType string `bson:"question_type" json:"question_type"`

	// admin add metadata of file and system will be found it
	Filename      string `bson:"filename" json:"filename"`
	AudioDuration string `bson:"audio_duration" json:"audio_duration"`
}

//go:generate mockery --name IQuizUseCase
type IQuizUseCase interface {
	FetchMany(ctx context.Context) ([]Quiz, error)
	FetchToDeleteMany(ctx context.Context) (*[]Quiz, error)
	UpdateOne(ctx context.Context, quizID string, quiz Quiz) error
	CreateOne(ctx context.Context, quiz *Quiz) error
	UpsertOne(c context.Context, id string, quiz *Quiz) (*Response, error)
	DeleteOne(ctx context.Context, quizID string) error
}
