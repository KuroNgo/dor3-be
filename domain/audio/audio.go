package audio_domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Audio struct {
	ID       int                `bson:"_id,omitempty" json:"id"`
	QuizID   primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	Filename string             `bson:"filename" json:"filename"`
	Format   string             `bson:"format" json:"format"`
	Path     string             `bson:"path" json:"path"`
}

type Response struct {
	QuizID   primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	Filename string             `bson:"filename" json:"filename"`
	Format   string             `bson:"format" json:"format"`
	Path     string             `bson:"path" json:"path"`
}

type IAudioRepository interface {
}
