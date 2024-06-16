package vocabulary_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Input struct {
	Id     primitive.ObjectID `bson:"_id" json:"_id"`
	UnitID primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Word          string `bson:"word" json:"word"`
	PartOfSpeech  string `bson:"part_of_speech" json:"part_of_speech"`
	Pronunciation string `bson:"pronunciation" json:"pronunciation"`
	Mean          string `bson:"mean" json:"mean"`
	ExampleVie    string `bson:"example_vie" json:"example_vie"`
	ExampleEng    string `bson:"example_eng" json:"example_eng"`
	ExplainVie    string `bson:"explain_vie" json:"explain_vie"`
	ExplainEng    string `bson:"explain_eng" json:"explain_eng"`
	FieldOfIT     string `bson:"field_of_it" json:"field_of_it"`
	LinkURL       string `bson:"link_url" json:"link_url"`
}

//go:generate mockery --name IVocabularyUseCase
type IVocabularyUseCase interface {
	FetchByLessonInBoth(ctx context.Context, unitName string) (SearchingResponse, error)
	FetchByWordInBoth(ctx context.Context, word string) (SearchingResponse, error)
	FetchManyInBoth(ctx context.Context, page string) (Response, error)
	UpdateIsFavouriteInUser(ctx context.Context, vocabularyID string, isFavourite int) error
	UpdateVocabularyProcess(ctx context.Context, vocabularyID string, process VocabularyProcess) error

	FindVocabularyIDByVocabularyConfigInAdmin(ctx context.Context, word string) (primitive.ObjectID, error)
	FetchByIdUnitInAdmin(ctx context.Context, idUnit string) ([]Vocabulary, error)
	GetAllVocabularyInAdmin(ctx context.Context) ([]string, error)
	GetVocabularyByIdInAdmin(ctx context.Context, id string) (Vocabulary, error)
	GetLatestVocabularyInAdmin(ctx context.Context) ([]string, error)

	CreateOneInAdmin(ctx context.Context, vocabulary *Vocabulary) error
	CreateOneByNameUnitInAdmin(ctx context.Context, vocabulary *Vocabulary) error
	UpdateOneInAdmin(ctx context.Context, vocabulary *Vocabulary) (*mongo.UpdateResult, error)
	UpdateOneImageInAdmin(ctx context.Context, vocabulary *Vocabulary) (*mongo.UpdateResult, error)
	UpdateOneAudioInAdmin(ctx context.Context, vocabulary *Vocabulary) error
	DeleteOneInAdmin(ctx context.Context, vocabularyID string) error
}
