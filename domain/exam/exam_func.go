package exam_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Input struct {
	LessonID      primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID        primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	VocabularyID  primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`
	Question      string             `bson:"question" json:"question"`
	Options       []string           `bson:"options" json:"options"`
	CorrectAnswer string             `bson:"correct_answer" json:"correct_answer"`
	Explanation   string             `bson:"explanation" json:"explanation"`
	Level         int                `bson:"level" json:"level"`
	QuestionType  string             `bson:"question_type" json:"question_type"`
}

type IExamUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByUnitID(ctx context.Context, unitID string) (Response, error)
	UpdateOne(ctx context.Context, examID string, exam Exam) error
	CreateOne(ctx context.Context, exam *Exam) error
	UpdateCompleted(ctx context.Context, examID string, isComplete int) error
	DeleteOne(ctx context.Context, examID string) error
}
