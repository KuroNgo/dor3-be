package exam_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID   primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Duration    string `bson:"duration" json:"duration"`
}

type IExamUseCase interface {
	FetchOneByUnitIDInUser(ctx context.Context, userID primitive.ObjectID, unitID string) (ExamProcessRes, error)
	UpdateCompletedInUser(ctx context.Context, userID primitive.ObjectID, exam *ExamProcess) error

	FetchManyInAdmin(ctx context.Context, page string) ([]Exam, DetailResponse, error)
	FetchExamByIDInAdmin(ctx context.Context, id string) (Exam, error)
	FetchManyByUnitIDInAdmin(ctx context.Context, unitID string, page string) ([]Exam, DetailResponse, error)
	FetchOneByUnitIDInAdmin(ctx context.Context, unitID string) (Exam, error)

	CreateOneInAdmin(ctx context.Context, exam *Exam) error
	UpdateOneInAdmin(ctx context.Context, exam *Exam) (*mongo.UpdateResult, error)
	DeleteOneInAdmin(ctx context.Context, examID string) error
}
