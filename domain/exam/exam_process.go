package exam_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionExamProcess = "exam_process"
)

type ExamProcess struct {
	ExamID       primitive.ObjectID `json:"exam_id" bson:"exam_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete   int                `json:"is_complete" bson:"is_complete"`
	JadeClaimNum int16              `json:"jade_claim" bson:"jade_claim"`
	IsClaimed    int                `json:"is_claimed" bson:"is_claimed"`
}

type ExamProcessRes struct {
	Exam         Exam               `json:"exam" bson:"exam"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete   int                `json:"is_complete" bson:"is_complete"`
	JadeClaimNum int16              `json:"jade_claim" bson:"jade_claim"`
	IsClaimed    int                `json:"is_claimed" bson:"is_claimed"`
}
