package exercise_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionExerciseProcess = "exercise_process"
)

type ExerciseProcess struct {
	ExerciseID   primitive.ObjectID `json:"exercise_id" bson:"exercise_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete   int                `json:"is_complete" bson:"is_complete"`
	JadeClaimNum int16              `json:"jade_claim" bson:"jade_claim"`
	IsClaimed    int                `json:"is_claimed" bson:"is_claimed"`
}

type ExerciseProcessRes struct {
	Exercise     Exercise           `json:"exercise" bson:"exercise"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete   int                `json:"is_complete" bson:"is_complete"`
	JadeClaimNum int16              `json:"jade_claim" bson:"jade_claim"`
	IsClaimed    int                `json:"is_claimed" bson:"is_claimed"`
}
