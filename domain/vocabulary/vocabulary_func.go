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
	LessonID      primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
}

type FetchByWordInput struct {
	Word string `bson:"word" json:"word"`
}

type FetchByLessonInput struct {
	Lesson string `bson:"lesson_id" json:"lesson_id"`
}

//go:generate mockery --name IVocabularyUseCase
type IVocabularyUseCase interface {
	FetchByID(ctx context.Context, vocabularyID string) (*Vocabulary, error)
	FetchMany(ctx context.Context) ([]Vocabulary, error)
	FetchByWord(ctx context.Context, word string) ([]Vocabulary, error)
	FetchByLesson(ctx context.Context, lessonName string) ([]Vocabulary, error)
	FetchToDeleteMany(ctx context.Context) (*[]Vocabulary, error)
	UpdateOne(ctx context.Context, vocabularyID string, vocabulary Vocabulary) error
	CreateOne(ctx context.Context, vocabulary *Vocabulary) error
	UpsertOne(c context.Context, id string, vocabulary *Vocabulary) (*Vocabulary, error)
	DeleteOne(ctx context.Context, vocabularyID string) error
}
