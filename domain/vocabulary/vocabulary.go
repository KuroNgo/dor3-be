package vocabulary_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CollectionVocabulary = "vocabulary"
)

type Vocabulary struct {
	Id            primitive.ObjectID `bson:"_id" json:"_id"`
	UnitID        primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	Word          string             `bson:"word" json:"word"`
	PartOfSpeech  string             `bson:"part_of_speech" json:"part_of_speech"`
	Mean          string             `bson:"mean" json:"mean"`
	Pronunciation string             `bson:"pronunciation" json:"pronunciation"`
	ExampleVie    string             `bson:"example_vie" json:"example_vie"`
	ExampleEng    string             `bson:"example_eng" json:"example_eng"`
	ExplainVie    string             `bson:"explain_vie" json:"explain_vie"`
	ExplainEng    string             `bson:"explain_eng" json:"explain_eng"`
	FieldOfIT     string             `bson:"field_of_it" json:"field_of_it"`
	LinkURL       string             `bson:"link_url" json:"link_url"`
	ImageURL      string             `bson:"image_url" json:"image_url"`
	IsFavourite   int                `bson:"is_favourite" json:"is_favourite"`
	WhoUpdates    string             `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	Vocabulary []Vocabulary
}

//go:generate mockery --name IVocabularyRepository
type IVocabularyRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FindUnitIDByUnitLevel(ctx context.Context, unitLevel int) (primitive.ObjectID, error)
	FetchByIdUnit(ctx context.Context, idUnit string) (Response, error)
	FetchByWord(ctx context.Context, word string) (Response, error)
	FetchByWord2(ctx context.Context, word string) (Response, error)
	FetchByLesson(ctx context.Context, unitName string) (Response, error)

	GetAllVocabulary(ctx context.Context) ([]string, error)
	GetLatestVocabulary(ctx context.Context) ([]string, error)

	CreateOne(ctx context.Context, vocabulary *Vocabulary) error
	CreateOneByNameUnit(ctx context.Context, vocabulary *Vocabulary) error

	UpdateOne(ctx context.Context, vocabulary *Vocabulary) (*mongo.UpdateResult, error)
	UpdateIsFavourite(ctx context.Context, vocabularyID string, isFavourite int) error
	UpdateOneAudio(c context.Context, vocabularyID string, linkURL string) error

	DeleteOne(ctx context.Context, vocabularyID string) error
}
