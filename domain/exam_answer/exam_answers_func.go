package exam_answer_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Input struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Answer string `bson:"answer" json:"answer"`
}

type IExamAnswerUseCase interface {
	FetchManyAnswerByQuestionIDInUser(ctx context.Context, questionID string, userID string) (Response, error)
	CreateOneInUser(ctx context.Context, examAnswer *ExamAnswer) error
	DeleteOneInUser(ctx context.Context, examID string) error
	DeleteAllAnswerByExamIDInUser(ctx context.Context, examID string) error
}
