package lesson_management

import (
	course_domain "clean-architecture/domain/course"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionCourseProcess = "unit_process"
)

type CourseProcess struct {
	CourseID   course_domain.Course `json:"course" bson:"course"`
	UserID     primitive.ObjectID   `json:"user_id" bson:"user_id"`
	IsComplete primitive.ObjectID   `json:"is_complete" bson:"is_complete"`
}
