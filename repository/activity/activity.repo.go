package activity_repository

import (
	activity_log_domain "clean-architecture/domain/activity_log"
	"clean-architecture/infrastructor/mongo"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type activityRepository struct {
	database           mongo.Database
	collectionActivity string
}

func (a *activityRepository) DeleteOneByTime(ctx context.Context, time time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (a *activityRepository) CreateOne(ctx context.Context, log activity_log_domain.ActivityLog) error {
	collectionActivity := a.database.Collection(a.collectionActivity)
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

func (a *activityRepository) FetchMany(ctx context.Context, page string) (activity_log_domain.Response, error) {
	collection := a.database.Collection(a.collectionActivity)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return activity_log_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return activity_log_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var activities []activity_log_domain.ActivityLog
	for cursor.Next(ctx) {
		var activity activity_log_domain.ActivityLog
		if err := cursor.Decode(&activity); err != nil {
			return activity_log_domain.Response{}, err
		}

		// Thêm activity vào slice activities
		activities = append(activities, activity)
	}

	activityRes := activity_log_domain.Response{
		ActivityLog: activities,
	}
	return activityRes, nil
}

func NewActivityRepository(db mongo.Database, collectionActivity string, collectionUser string) activity_log_domain.IActivityRepository {
	return &activityRepository{
		database:           db,
		collectionActivity: collectionActivity,
	}
}
