package course_domain

import "context"

type Input struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

//go:generate mockery --name ICourseUseCase
type ICourseUseCase interface {
	FetchMany(ctx context.Context) (Response, error)
	UpdateOne(ctx context.Context, courseID string, course Course) error
	CreateOne(ctx context.Context, course *Course) error
	DeleteOne(ctx context.Context, courseID string) error
}
