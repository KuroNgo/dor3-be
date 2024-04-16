package activity_repository

import (
	activity_log_domain "clean-architecture/domain/activity_log"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type activityRepository struct {
	database           *mongo.Database
	collectionActivity string
}

func NewActivityRepository(db *mongo.Database, collectionActivity string) activity_log_domain.IActivityRepository {
	return &activityRepository{
		database:           db,
		collectionActivity: collectionActivity,
	}
}

func (a *activityRepository) CreateOne(ctx context.Context, log activity_log_domain.ActivityLog) error {
	collectionActivity := a.database.Collection(a.collectionActivity)

	now := time.Now()
	tomorrow := now.Add(24 * 30 * time.Hour)
	expireTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.UTC)

	log.ExpireAt = expireTime

	_, err := collectionActivity.InsertOne(ctx, &log)
	if err != nil {
		return err
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
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var activities []activity_log_domain.ActivityLog
	for cursor.Next(ctx) {
		var activity activity_log_domain.ActivityLog
		if err := cursor.Decode(&activity); err != nil {
			return activity_log_domain.Response{}, err
		}

		activity.ActivityTime = activity.ActivityTime.Add(7 * time.Hour)

		// Thêm activity vào slice activities
		activities = append(activities, activity)
	}

	activityRes := activity_log_domain.Response{
		ActivityLog: activities,
	}
	return activityRes, nil
}
