package user_process

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserProcess struct {
	ID            primitive.ObjectID `bson:"id" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	CourseID      primitive.ObjectID `bson:"course_id" json:"course_id"`
	LessonID      primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	Score         int                `bson:"score" json:"score"`
	CompletedDate time.Time          `bson:"completed_date" json:"completed_date"`
}

type UserProcessInput struct {
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	CourseID      primitive.ObjectID `bson:"course_id" json:"course_id"`
	LessonID      primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	Score         int                `bson:"score" json:"score"`
	CompletedDate time.Time          `bson:"completed_date" json:"completed_date"`
}
