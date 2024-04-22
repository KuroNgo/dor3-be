package exam_answer_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Input struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`

	Content string `bson:"content" json:"content"`
}

type IExamAnswerUseCase interface {
	FetchManyAnswerByUserIDAndQuestionID(ctx context.Context, questionID string, userID string) (Response, error)
	CreateOne(ctx context.Context, examAnswer *ExamAnswer) error
	DeleteOne(ctx context.Context, examID string) error
	DeleteAllAnswerByExamID(ctx context.Context, examID string) error
}
