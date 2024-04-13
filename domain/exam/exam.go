package exam_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionExam = "exam"
)

type Exam struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	LessonID      primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID        primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	VocabularyID  primitive.ObjectID `bson:"vocabulary_id" json:"vocabulary_id"`
	Question      string             `bson:"question" json:"question"`
	Options       []string           `bson:"options" json:"options"`
	CorrectAnswer string             `bson:"correct_answer" json:"correct_answer"`
	Explanation   string             `bson:"explanation" json:"explanation"`
	QuestionType  string             `bson:"question_type" json:"question_type"`
	Level         int                `bson:"level" json:"level"`

	// admin add metadata of file and system will be found it
	Filename      string `bson:"filename" json:"filename"`
	AudioDuration string `bson:"audio_duration" json:"audio_duration"`

	IsComplete int       `bson:"is_complete" json:"is_complete"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	Exam  []Exam
	Count int64 `bson:"count" json:"count"`
}

type IExamRepository interface {
	FetchMany(ctx context.Context) (Response, error)
	FetchManyByUnitID(ctx context.Context, unitID string) (Response, error)
	UpdateOne(ctx context.Context, examID string, exam Exam) error
	CreateOne(ctx context.Context, exam *Exam) error
	UpdateCompleted(ctx context.Context, examID string, isComplete int) error
	DeleteOne(ctx context.Context, examID string) error
}
