package lesson_management_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionUnitProcess = "unit_process"
)

type UnitProcess struct {
	UnitID             primitive.ObjectID `json:"unit_id" bson:"unit_id"`
	UserID             primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete         int                `json:"is_complete" bson:"is_complete"`
	ExamIsComplete     int                `bson:"exam_is_complete" json:"exam_is_complete"`
	ExerciseIsComplete int                `bson:"exercise_is_complete" json:"exercise_is_complete"`
	QuizIsComplete     int                `bson:"quiz_is_complete" json:"quiz_is_complete"`
	TotalScore         int32              `json:"total_score" bson:"total_score"`
}
