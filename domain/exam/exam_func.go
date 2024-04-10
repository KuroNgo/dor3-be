package exam_domain

import "context"

type Input struct {
	Question      string   `bson:"question" json:"question"`
	Options       []string `bson:"options" json:"options"`
	CorrectAnswer string   `bson:"correct_answer" json:"correct_answer"`
	Explanation   string   `bson:"explanation" json:"explanation"`
	Level         int      `bson:"level" json:"level"`
	QuestionType  string   `bson:"question_type" json:"question_type"`
}

type IExamUseCase interface {
	FetchMany(ctx context.Context) (Response, error)
	UpdateOne(ctx context.Context, examID string, exam Exam) error
	CreateOne(ctx context.Context, exam *Exam) error
	UpsertOne(c context.Context, id string, exam *Exam) (Response, error)
	DeleteOne(ctx context.Context, examID string) error
}