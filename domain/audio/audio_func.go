package audio_domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AutoMatch struct {
	QuizID primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`

	// admin add file and system will be found it
	Filename      string `bson:"filename" json:"filename"`
	AudioDuration string `bson:"audio_duration" json:"audio_duration"`
}

//go:generate mockery --name IAudioUseCase
type IAudioUseCase interface {
	FetchMany(ctx context.Context) ([]Audio, error)
	FetchToDeleteMany(ctx context.Context) (*[]Audio, error)
	UpdateOne(ctx context.Context, audioID string, audio Audio) error

	// CreateOne needn't input, because the system will be found information file
	CreateOne(ctx context.Context, audio *Audio) error

	DeleteOne(ctx context.Context, audioID string) error
	DeleteMany(ctx context.Context, audioID ...string) error
}
