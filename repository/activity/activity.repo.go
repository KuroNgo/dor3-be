package activity_repository

import (
	activity_log_domain "clean-architecture/domain/activity_log"
	"clean-architecture/infrastructor/mongo"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type activityRepository struct {
	database           mongo.Database
	collectionActivity string
	collectionUser     string
}

func (a *activityRepository) DeleteOneByTime(ctx context.Context, time time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (a *activityRepository) CreateOne(ctx context.Context, log activity_log_domain.ActivityLog) error {
	// Kiểm tra log.UserID
	if log.UserID == primitive.NilObjectID {
		return errors.New("userID is empty")
	}

	errChan := make(chan error, 2)

	go func() {
		userFilter := bson.M{"_id": log.UserID}
		userCount, err := a.database.Collection(a.collectionUser).CountDocuments(ctx, userFilter)
		if err != nil {
			errChan <- fmt.Errorf("error checking user existence: %v", err)
			return
		}
		if userCount == 0 {
			errChan <- errors.New("user does not exist")
			return
		}
		errChan <- nil
	}()

	go func() {
		_, err := a.database.Collection(a.collectionActivity).InsertOne(ctx, log)
		if err != nil {
			errChan <- fmt.Errorf("error inserting activity: %v", err)
			return
		}
		errChan <- nil
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

func (a *activityRepository) DeleteOne(ctx context.Context, logID string) error {
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
