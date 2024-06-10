package lesson_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionLessonProcess = "lesson_process"
)

type LessonProcess struct {
	LessonID     primitive.ObjectID `json:"lesson_id" bson:"lesson_id"`
	CourseID     primitive.ObjectID `json:"course_id" bson:"course_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete   int                `json:"is_complete" bson:"is_complete"`
	UnitComplete []int              `bson:"unit_complete" json:"unit_complete"`
	TotalScore   int32              `json:"total_score" bson:"total_score"`
}

type LessonProcessResponse struct {
	Lesson       Lesson             `json:"lesson" bson:"lesson"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete   int                `json:"is_complete" bson:"is_complete"`
	UnitComplete []int              `bson:"unit_complete" json:"unit_complete"`
	TotalScore   int32              `json:"total_score" bson:"total_score"`
}

type LessonProcessResponseList []LessonProcessResponse

func (l LessonProcessResponseList) Len() int {
	return len(l)
}

func (l LessonProcessResponseList) Less(i, j int) bool {
	return l[i].Lesson.Level < l[j].Lesson.Level
}

func (l LessonProcessResponseList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
