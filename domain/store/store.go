package store

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Store struct {
	StoreID  primitive.ObjectID `bson:"_id" json:"-"`
	Name     string             `bson:"name" json:"name"`
	Avatar   string             `bson:"avatar" json:"avatar"`
	Bio      string             `bson:"bio" json:"bio"`
	OpenAt   time.Time          `bson:"open_at" json:"open_at"`
	ClosedAt time.Time          `bson:"closed_at" json:"closed_at"`
	Phone    string             `bson:"phone" json:"phone"`
	Location string             `bson:"location" json:"location"`
}
