package exam_management_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionExamProcess = "exam_process"
)

type ExamProcess struct {
	ExamID     primitive.ObjectID `json:"exam_id" bson:"exam_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete int                `json:"is_complete" bson:"is_complete"`
}
