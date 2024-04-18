package exam_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Input struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID     primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID       primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	VocabularyID primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`

	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Duration    time.Duration `bson:"duration" json:"duration"`
}

type IExamUseCase interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchManyByUnitID(ctx context.Context, unitID string) (Response, error)

	CreateOne(ctx context.Context, exam *Exam) error
	UpdateOne(ctx context.Context, exam *Exam) (*mongo.UpdateResult, error)
	UpdateCompleted(ctx context.Context, examID string, isComplete int) error
	DeleteOne(ctx context.Context, examID string) error
}
