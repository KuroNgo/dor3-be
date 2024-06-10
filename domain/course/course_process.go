package course_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionCourseProcess = "course_process"
)

type CourseProcess struct {
	CourseID   primitive.ObjectID `json:"course_id" bson:"course_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	IsComplete int                `json:"is_complete" bson:"is_complete"`
}
