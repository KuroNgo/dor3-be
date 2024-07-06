package quiz_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionQuizProcess = "quiz_process"
)

type QuizProcess struct {
	QuizID       primitive.ObjectID `json:"quiz_id" bson:"quiz_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete   int                `json:"is_complete" bson:"is_complete"`
	JadeClaimNum int16              `json:"jade_claim" bson:"jade_claim"`
	IsClaimed    int                `json:"is_claimed" bson:"is_claimed"`
}

type QuizProcessRes struct {
	Quiz         Quiz               `json:"quiz" bson:"quiz"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete   int                `json:"is_complete" bson:"is_complete"`
	JadeClaimNum int16              `json:"jade_claim" bson:"jade_claim"`
	IsClaimed    int                `json:"is_claimed" bson:"is_claimed"`
}
