package quiz_domain

import "context"

type Input struct {
	Question      string   `bson:"question" json:"question"`
	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`
	Explanation   string   `bson:"explanation" json:"explanation"`

	// QuestionType can be included checkbox, check radius or write correct answer
	QuestionType string `bson:"question_type" json:"question_type"`
}

//go:generate mockery --name IQuizUseCase
type IQuizUseCase interface {
	Fetch(ctx context.Context) ([]Quiz, error)
	FetchToDelete(ctx context.Context) (*[]Quiz, error)
	Update(ctx context.Context, quizID string, quiz Quiz) error
	Create(ctx context.Context, quiz *Input) error
	Upsert(c context.Context, question string, quiz *Quiz) (*Response, error)
	Delete(ctx context.Context, quizID string) error
}
