package exam_question_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionExamQuestionProcess = "exam_question_process"
)

type ExamQuestionProcess struct {
	QuestionID primitive.ObjectID `json:"exam_id" bson:"exam_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete int                `json:"is_complete" bson:"is_complete"`
	IsTrue     int                `json:"is_true" bson:"is_true"`
}

type ExamProcessRes struct {
	Question   ExamQuestion       `json:"exam_question" bson:"exam_question"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete int                `json:"is_complete" bson:"is_complete"`
	IsTrue     int                `json:"is_true" bson:"is_true"`
}
