package exercise_options

import "go.mongodb.org/mongo-driver/mongo"

type exerciseOptionsRepository struct {
	database           mongo.Database
	collectionQuestion string
	collectionOptions  string
}
