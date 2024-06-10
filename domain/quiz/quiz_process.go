package quiz_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionQuizProcess = "quiz_process"
)

type QuizProcess struct {
	QuizID     primitive.ObjectID `json:"quiz_id" bson:"quiz_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete int                `json:"is_complete" bson:"is_complete"`
}
