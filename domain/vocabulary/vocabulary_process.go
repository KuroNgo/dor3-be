package vocabulary_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionVocabularyProcess = "vocabulary_process"
)

type VocabularyProcess struct {
	VocabularyID primitive.ObjectID `json:"vocabulary_id" bson:"vocabulary_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsFavourite  int                `json:"is_favourite" bson:"is_favourite"`
}
