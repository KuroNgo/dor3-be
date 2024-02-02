package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Lesson struct {
	ID       primitive.ObjectID `bson:"id" json:"id"`
	CourseID primitive.ObjectID `bson:"course_id" json:"course_id"`
	Name     string             `bson:"name" json:"name"`
	Content  string             `bson:"content" json:"content"`
}

type LessonInput struct {
	// default Course = 1
	CourseID primitive.ObjectID `bson:"course_id" json:"course_id"`
	Name     string             `bson:"name" json:"name"`
	Content  string             `bson:"content" json:"content"`
}
