package lesson_management

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionVocabularyFavourite = "vocabulary_favourite"
)

type VocabularyProcess struct {
	VocabularyID primitive.ObjectID `json:"vocabulary_id" bson:"vocabulary_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsFavourite  int                `json:"is_favourite" bson:"is_favourite"`
}
