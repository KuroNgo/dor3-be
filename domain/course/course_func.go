package course_domain

import "context"

type Input struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Level       int    `bson:"level" json:"level"`
}

type ICourseUseCase interface {
	FetchByID(ctx context.Context, courseID string) (*Response, error)
	FetchMany(ctx context.Context) ([]Course, error)
	FetchToDeleteMany(ctx context.Context) (*[]Course, error)
	UpdateOne(ctx context.Context, courseID string, course Course) error
	CreateOne(ctx context.Context, course *Course) error
	UpsertOne(ctx context.Context, id string, course *Course) (*Response, error)
	DeleteOne(ctx context.Context, courseID string) error
}
