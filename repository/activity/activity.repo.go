package activity_repository

import (
	activity_log_domain "clean-architecture/domain/activity_log"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
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

	cal := make(chan int64)

	go func() {
		count, err := collection.CountDocuments(ctx, bson.D{})
		if err != nil {
			return
		}

		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			cal <- cal1
		}
	}()

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

	// Sắp xếp slice activities theo thời gian giảm dần (từ mới nhất đến cũ nhất)
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].ActivityTime.After(activities[j].ActivityTime)
	})

	count := <-cal
	activityRes := activity_log_domain.Response{
		Page:        count,
		PageCurrent: int64(pageNumber),
		ActivityLog: activities,
	}
	return activityRes, nil
}
