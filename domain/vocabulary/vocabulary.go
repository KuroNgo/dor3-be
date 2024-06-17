package vocabulary_domain

import (
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
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
	VideoURL      string `bson:"video_url" json:"video_url"`
	ImageURL      string `bson:"image_url" json:"image_url"`
	AssetURL      string `bson:"asset_url" json:"asset_url"`

	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	WhoUpdates string    `bson:"who_updates" json:"who_updates"`
}

type VocabularyResponse struct {
	Vocabulary Vocabulary           `bson:"vocabulary" json:"vocabulary"`
	Unit       unit_domain.Unit     `bson:"unit_id" json:"unit_id"`
	Lesson     lesson_domain.Lesson `bson:"lesson" json:"lesson"`
}

type Response struct {
	Page               int64                `bson:"page" json:"page"`
	CurrentPage        int                  `bson:"current_page" json:"current_page"`
	VocabularyResponse []VocabularyResponse `bson:"vocabulary" json:"vocabulary"`
}

type SearchingResponse struct {
	CountVocabularySearch int64        `bson:"count_vocabulary_search" json:"count_vocabulary_search"`
	Vocabulary            []Vocabulary `bson:"vocabulary" json:"vocabulary"`
}

//go:generate mockery --name IVocabularyRepository
type IVocabularyRepository interface {
	FetchManyInBoth(ctx context.Context, page string) (Response, error)
	FetchByWordInBoth(ctx context.Context, word string) (SearchingResponse, error)
	FetchByLessonInBoth(ctx context.Context, unitName string) (SearchingResponse, error)
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
