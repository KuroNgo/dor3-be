package activity_repository

import (
	activity_log_domain "clean-architecture/domain/activity_log"
	"clean-architecture/infrastructor/mongo"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type activityRepository struct {
	database           mongo.Database
	collectionActivity string
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
