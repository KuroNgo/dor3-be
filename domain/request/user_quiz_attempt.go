package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserQuizAttempt struct {
	ID     primitive.ObjectID `bson:"id" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuizID primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	Answer int                `bson:"answer" json:"answer"`
	Score  int                `bson:"score" json:"score"`
}
