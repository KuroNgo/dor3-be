package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	ID          primitive.ObjectID `bson:"id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Level       int                `bson:"level" json:"level"`
}

type CourseInput struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Level       int    `bson:"level" json:"level"`
}
