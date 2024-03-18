package vocabulary_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Input struct {
	Word          string             `bson:"word" json:"word"`
	PartOfSpeech  string             `bson:"part_of_speech" json:"part_of_speech"`
	Pronunciation string             `bson:"pronunciation" json:"pronunciation"`
	Example       string             `bson:"example" json:"example"`
	FieldOfIT     string             `bson:"field_of_it" json:"field_of_it"`
	LinkURL       string             `bson:"link_url" json:"link_url"`
	UnitID        primitive.ObjectID `bson:"unit_id" json:"unit_id"`
}

type FetchByWordInput struct {
	Word string `bson:"word" json:"word"`
}

type FetchByLessonInput struct {
	Lesson string `bson:"lesson_id" json:"lesson_id"`
}

//go:generate mockery --name IVocabularyUseCase
type IVocabularyUseCase interface {
	FetchMany(ctx context.Context) ([]Response, error)
	FetchByWord(ctx context.Context, word string) ([]Response, error)
	FetchByLesson(ctx context.Context, lessonName string) ([]Response, error)
	UpdateOne(ctx context.Context, vocabularyID string, vocabulary Vocabulary) error
	CreateOne(ctx context.Context, vocabulary *Vocabulary) error
	UpsertOne(c context.Context, id string, vocabulary *Vocabulary) (*Response, error)
	DeleteOne(ctx context.Context, vocabularyID string) error
}
