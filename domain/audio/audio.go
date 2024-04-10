package audio_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionAudio = "audio"
)

type Audio struct {
	Id primitive.ObjectID `bson:"_id" json:"_id"`
	//QuizID primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`

	// admin add metadata of file and system will be found it
	Filename string `bson:"filename" json:"filename"`
	Size     int64  `bson:"size" json:"size"`
}

type Response struct {
	Audio []Audio
	Count int64 `json:"count"`
}

//go:generate mockery --name IAudioRepository
type IAudioRepository interface {
	FetchMany(ctx context.Context) (Response, error)
	UpdateOne(ctx context.Context, audioID string, audio Audio) error
	CreateOne(ctx context.Context, audio *Audio) error
	DeleteOne(ctx context.Context, audioID string) error
	DeleteMany(ctx context.Context, audioID ...string) error
}
