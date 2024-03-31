package activity_repository

import (
	activity_log_domain "clean-architecture/domain/activity_log"
	"clean-architecture/infrastructor/mongo"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type activityRepository struct {
	database           mongo.Database
	collectionActivity string
	collectionUser     string
}

func (a *activityRepository) CreateOne(ctx context.Context, log activity_log_domain.ActivityLog) error {
	collectionActivity := a.database.Collection(a.collectionActivity)
	collectionUser := a.database.Collection(a.collectionUser)

	filterUser := bson.M{"_id": log.UserID}

	countUser, err := collectionUser.CountDocuments(ctx, filterUser)
	if err != nil {
		return err
	}
	if countUser <= 0 {
		return errors.New("userId do not exist")
	}

	_, err = collectionActivity.InsertOne(ctx, log)
	return err
}

func (a *activityRepository) DeleteOne(ctx context.Context) error {
	collectionActivity := a.database.Collection(a.collectionActivity)

	timeThreshold := time.Now().AddDate(0, 0, -30)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"activity_time": bson.M{"$lte": timeThreshold},
			},
		},
	}

	cursor, err := collectionActivity.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var activity bson.M
		if err := cursor.Decode(&activity); err != nil {
			return err
		}

		id, ok := activity["_id"].(primitive.ObjectID)
		if !ok {
			return errors.New("failed to parse activity ID")
		}

		_, err := collectionActivity.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *activityRepository) FetchMany(ctx context.Context) ([]activity_log_domain.ActivityLog, error) {
	//TODO implement me
	panic("implement me")
}

func (a *activityRepository) FetchByUserName(ctx context.Context, username string) (activity_log_domain.ActivityLog, error) {
	//TODO implement me
	panic("implement me")
}

func NewActivityRepository(db mongo.Database, collectionActivity string, collectionUser string) activity_log_domain.IActivityUseCase {
	return &activityRepository{
		database:           db,
		collectionActivity: collectionActivity,
		collectionUser:     collectionUser,
	}
}
