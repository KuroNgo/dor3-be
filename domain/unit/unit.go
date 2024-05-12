package unit_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionUnit = "unit"
)

type Unit struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`

	Name     string `bson:"name" json:"name"`
	ImageURL string `bson:"image_url" json:"image_url"`
	Level    int    `bson:"level" json:"level"`

	IsComplete int       `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoCreate  string    `bson:"who_create" json:"who_create"`
	Learner    string    `bson:"learner" json:"learner"`
}

type UnitResponse struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`

	Name     string `bson:"name" json:"name"`
	ImageURL string `bson:"image_url" json:"image_url"`
	Level    int    `bson:"level" json:"level"`

	IsComplete int       `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoCreate  string    `bson:"who_create" json:"who_create"`
	Learner    string    `bson:"learner" json:"learner"`

	ExamIsComplete     int   `bson:"exam_is_complete" json:"exam_is_complete"`
	ExerciseIsComplete int   `bson:"exercise_is_complete" json:"exercise_is_complete"`
	QuizIsComplete     int   `bson:"quiz_is_complete" json:"quiz_is_complete"`
	CountVocabulary    int32 `json:"count_vocabulary"`
}

type DetailResponse struct {
	CountUnit   int64 `json:"count_unit"`
	Page        int64 `json:"page"`
	CurrentPage int   `json:"current_page"`
}

//go:generate mockery --name IUnitRepository
type IUnitRepository interface {
	FetchMany(ctx context.Context, page string) ([]UnitResponse, DetailResponse, error)
	FindLessonIDByLessonName(ctx context.Context, lessonName string) (primitive.ObjectID, error)
	FetchByIdLesson(ctx context.Context, idLesson string, page string) ([]UnitResponse, DetailResponse, error)

	CreateOne(ctx context.Context, unit *Unit) error
	CreateOneByNameLesson(ctx context.Context, unit *Unit) error

	DeleteOne(ctx context.Context, unitID string) error

	// UpdateComplete automation
	UpdateComplete(ctx context.Context, update *Unit) error
	UpdateOne(ctx context.Context, unit *Unit) (*mongo.UpdateResult, error)
	CheckLessonComplete(ctx context.Context, lessonID primitive.ObjectID) (bool, error)
}
