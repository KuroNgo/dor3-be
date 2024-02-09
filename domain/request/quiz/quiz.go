package quiz

import "go.mongodb.org/mongo-driver/bson/primitive"

type Quiz struct {
	ID       primitive.ObjectID
	LessonID primitive.ObjectID
}
