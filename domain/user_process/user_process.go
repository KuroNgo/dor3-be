package user_process

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserProcess struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	CourseID      primitive.ObjectID `bson:"course_id" json:"course_id"`
	LessonID      primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID        primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	Score         int                `bson:"score" json:"score"`
	ProcessStatus int                `bson:"process_status" json:"process_status"`
	CompletedDate time.Time          `bson:"completed_date" json:"completed_date"`
}

type UserProcessAuto struct {
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	CourseID      primitive.ObjectID `bson:"course_id" json:"course_id"`
	LessonID      primitive.ObjectID `bson:"lesson_id" json:"lesson_id"`
	UnitID        primitive.ObjectID `bson:"unit_id" json:"unit_id"`
	Score         int                `bson:"score" json:"score"`
	ProcessStatus int                `bson:"process_status" json:"process_status"`
	CompletedDate time.Time          `bson:"completed_date" json:"completed_date"`
}

type IUserProcessRepository interface {
}
