package activity_repository

import (
	activity_log_domain "clean-architecture/domain/activity_log"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type activityRepositoryV2 struct {
	database           mongo.Database
	collectionActivity string
}

func NewActivityRepositoryV2(db mongo.Database, collectionName string) activity_log_domain.IActivityRepositoryV2 {
	return &activityRepositoryV2{
		database:           db,
		collectionActivity: collectionName,
	}
}

func (a *activityRepositoryV2) CreateOne(ctx context.Context, log activity_log_domain.ActivityLog) error {
	collectionActivity := a.database.Collection(a.collectionActivity)

	// TTL index
	ttlIndex := mongo.IndexModel{
		Keys:    bson.M{"activity_time": 1},                               // Chỉ định trường và hướng sắp xếp
		Options: options.Index().SetExpireAfterSeconds(30 * 24 * 60 * 60), // Thời gian sống của bản ghi (30 ngày)
	}
	_, err := collectionActivity.Indexes().CreateOne(ctx, ttlIndex)
	if err != nil {
		return err
	}

	filter := bson.M{"activity_time": log.ActivityTime}

	// check exists with CountDocuments
	count, err := collectionActivity.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("warning")
	}
	_, err = collectionActivity.InsertOne(ctx, log)
	return nil
}
