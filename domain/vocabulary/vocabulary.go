package vocabulary_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	CollectionVocabulary = "vocabulary"
)

type Vocabulary struct {
	Id     primitive.ObjectID `bson:"_id" json:"_id"`
	UnitID primitive.ObjectID `bson:"unit_id" json:"unit_id"`

	Word          string `bson:"word" json:"word"`
	WordForConfig string `bson:"word_for_config" json:"word_for_config"`
	PartOfSpeech  string `bson:"part_of_speech" json:"part_of_speech"`
	Mean          string `bson:"mean" json:"mean"`
	Pronunciation string `bson:"pronunciation" json:"pronunciation"`
	ExampleVie    string `bson:"example_vie" json:"example_vie"`
	ExampleEng    string `bson:"example_eng" json:"example_eng"`
	ExplainVie    string `bson:"explain_vie" json:"explain_vie"`
	ExplainEng    string `bson:"explain_eng" json:"explain_eng"`
	FieldOfIT     string `bson:"field_of_it" json:"field_of_it"`
	LinkURL       string `bson:"link_url" json:"link_url"`
	ImageURL      string `bson:"image_url" json:"image_url"`

	IsFavourite int       `bson:"is_favourite" json:"is_favourite"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates  string    `bson:"who_updates" json:"who_updates"`
}

type Response struct {
	Page        int64        `bson:"page" json:"page"`
	CurrentPage int          `bson:"current_page" json:"current_page"`
	Vocabulary  []Vocabulary `bson:"vocabulary" json:"vocabulary"`
}

type SearchingResponse struct {
	CountVocabularySearch int64        `bson:"count_vocabulary_search" json:"count_vocabulary_search"`
	Vocabulary            []Vocabulary `bson:"vocabulary" json:"vocabulary"`
}

//go:generate mockery --name IVocabularyRepository
type IVocabularyRepository interface {
	FetchMany(ctx context.Context, page string) (Response, error)
	FindUnitIDByUnitLevel(ctx context.Context, unitLevel int) (primitive.ObjectID, error)
	FindVocabularyIDByVocabularyConfig(ctx context.Context, word string) (primitive.ObjectID, error)
	FetchByIdUnit(ctx context.Context, idUnit string) (Response, error)
	FetchByWord(ctx context.Context, word string) (SearchingResponse, error)
	FetchByLesson(ctx context.Context, unitName string) (SearchingResponse, error)

	GetAllVocabulary(ctx context.Context) ([]string, error)
	GetVocabularyById(ctx context.Context, id string) (Vocabulary, error)
	GetLatestVocabulary(ctx context.Context) ([]string, error)

	CreateOne(ctx context.Context, vocabulary *Vocabulary) error
	CreateOneByNameUnit(ctx context.Context, vocabulary *Vocabulary) error

	UpdateOne(ctx context.Context, vocabulary *Vocabulary) (*mongo.UpdateResult, error)
	UpdateOneImage(ctx context.Context, vocabulary *Vocabulary) (*mongo.UpdateResult, error)
	UpdateIsFavourite(ctx context.Context, vocabularyID string, isFavourite int) error
	UpdateOneAudio(ctx context.Context, vocabulary *Vocabulary) error

	DeleteOne(ctx context.Context, vocabularyID string) error
}
