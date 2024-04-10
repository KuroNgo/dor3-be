package vocabulary_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionVocabulary = "vocabulary"
)

type Vocabulary struct {
	Id            primitive.ObjectID `bson:"_id" json:"_id"`
	UnitID        primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	Word          string             `bson:"word" json:"word"`
	PartOfSpeech  string             `bson:"part_of_speech" json:"part_of_speech"`
	Pronunciation string             `bson:"pronunciation" json:"pronunciation"`
	ExampleVie    string             `bson:"example_vie" json:"example_vie"`
	ExampleEng    string             `bson:"example_eng" json:"example_eng"`
	ExplainVie    string             `bson:"explain_vie" json:"explain_vie"`
	ExplainEng    string             `bson:"explain_eng" json:"explain_eng"`
	FieldOfIT     string             `bson:"field_of_it" json:"field_of_it"`
	LinkURL       string             `bson:"link_url" json:"link_url"`
}

type Response struct {
	Vocabulary []Vocabulary
}

//go:generate mockery --name IVocabularyRepository
type IVocabularyRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FindUnitIDByUnitName(ctx context.Context, unitName string) (primitive.ObjectID, error)
	GetAllVocabulary(ctx context.Context) ([]string, error)
	FetchByIdUnit(ctx context.Context, idUnit string) (Response, error)
	FetchByWord(ctx context.Context, word string) (Response, error)
	FetchByLesson(ctx context.Context, unitName string) (Response, error)
	UpdateOne(ctx context.Context, vocabularyID string, vocabulary Vocabulary) error
	CreateOne(ctx context.Context, vocabulary *Vocabulary) error
	CreateOneByNameUnit(ctx context.Context, vocabulary *Vocabulary) error
	UpsertOne(c context.Context, id string, vocabulary *Vocabulary) (Response, error)
	UpdateOneAudio(c context.Context, vocabularyID string, linkURL string) error
	DeleteOne(ctx context.Context, vocabularyID string) error
}
