package quiz_domain

import "context"

type Input struct {
	Question      string   `bson:"question" json:"question"`
	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`
	Explanation   string   `bson:"explanation" json:"explanation"`
	Level         int      `bson:"level" json:"level"`

	// QuestionType can be included checkbox, check radius or write correct answer
	QuestionType string `bson:"question_type" json:"question_type"`
	Skill        string `bson:"skill" json:"skill"`
}

//go:generate mockery --name IQuizUseCase
type IQuizUseCase interface {
	FetchMany(ctx context.Context) (Response, error)
	//FetchTenQuizButEnoughAllSkill(ctx context.Context) ([]Response, error)
	UpdateOne(ctx context.Context, quizID string, quiz Quiz) error
	CreateOne(ctx context.Context, quiz *Quiz) error
	DeleteOne(ctx context.Context, quizID string) error
}
