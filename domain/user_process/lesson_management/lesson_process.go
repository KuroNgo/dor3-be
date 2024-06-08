package lesson_management

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionLessonProcess = "unit_process"
)

type LessonProcess struct {
	//Lesson lesson_domain.Lesson `json:"lesson" bson:"lesson"`
	LessonID     primitive.ObjectID `json:"lesson_id" bson:"lesson_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete   int                `json:"is_complete" bson:"is_complete"`
	UnitComplete []int              `bson:"unit_complete" json:"unit_complete"`
	TotalScore   int32              `json:"total_score" bson:"total_score"`
}
