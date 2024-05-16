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
	FetchMany(ctx context.Context, page string) (Response, error)
	FetchByIdUnit(ctx context.Context, idUnit string) (Response, error)
	FindUnitIDByUnitLevel(ctx context.Context, unitLevel int, fieldOfIT string) (primitive.ObjectID, error)
	FindVocabularyIDByVocabularyName(ctx context.Context, word string) (primitive.ObjectID, error)
	FetchByWord(ctx context.Context, word string) (SearchingResponse, error)
	FetchByLesson(ctx context.Context, unitName string) (SearchingResponse, error)

	GetAllVocabulary(ctx context.Context) ([]string, error)
	GetVocabularyById(ctx context.Context, id string) (Vocabulary, error)
	GetLatestVocabulary(ctx context.Context) ([]string, error)

	UpdateOne(ctx context.Context, vocabulary *Vocabulary) (*mongo.UpdateResult, error)
	UpdateOneAudio(ctx context.Context, vocabulary *Vocabulary) error
	UpdateIsFavourite(ctx context.Context, vocabularyID string, isFavourite int) error
	UpdateOneImage(ctx context.Context, vocabulary *Vocabulary) (*mongo.UpdateResult, error)

	CreateOne(ctx context.Context, vocabulary *Vocabulary) error
	CreateOneByNameUnit(ctx context.Context, vocabulary *Vocabulary) error

	DeleteOne(ctx context.Context, vocabularyID string) error
}
