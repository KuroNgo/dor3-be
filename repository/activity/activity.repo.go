package activity_repository

import (
	activity_log_domain "clean-architecture/domain/activity_log"
	admin_domain "clean-architecture/domain/admin"
	"clean-architecture/internal"
	"clean-architecture/internal/cache/memory"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"sync"
	"time"
)

type activityRepository struct {
	database           *mongo.Database
	collectionActivity string
	collectionAdmin    string
}

func NewActivityRepository(db *mongo.Database, collectionActivity string, collectionAdmin string) activity_log_domain.IActivityRepository {
	return &activityRepository{
		database:           db,
		collectionActivity: collectionActivity,
		collectionAdmin:    collectionAdmin,
	}
}

var (
	statisticsCache = memory.NewTTL[string, activity_log_domain.Statistics]()
	wg              sync.WaitGroup
)

const (
	cacheTTL = 24 * time.Hour
)

func (a *activityRepository) CreateOne(ctx context.Context, log activity_log_domain.ActivityLog) error {
	collectionActivity := a.database.Collection(a.collectionActivity)

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	expireTime := time.Date(tomorrow.Year(), tomorrow.Month()+1, tomorrow.Day(), 0, 0, 0, 0, time.UTC)

	log.ExpireAt = expireTime

	_, err := collectionActivity.InsertOne(ctx, &log)
	if err != nil {
		return err
	}

	// Tạo TTL Index
	index := mongo.IndexModel{
		Keys:    bson.M{"expire_at": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, err = collectionActivity.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		statisticsCache.Clear()
	}()
	wg.Wait()

	return nil
}

func (a *activityRepository) FetchMany(ctx context.Context, page string) (activity_log_domain.Response, error) {
	errCh := make(chan error, 1)

	collection := a.database.Collection(a.collectionActivity)
	collectionAdmin := a.database.Collection(a.collectionAdmin)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return activity_log_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip)).SetSort(bson.D{{"_id", -1}})

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return activity_log_domain.Response{}, err
	}

	totalPages := (count + int64(perPage) - 1) / int64(perPage)

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return activity_log_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var activities []activity_log_domain.ActivityLog
	activities = make([]activity_log_domain.ActivityLog, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var activity activity_log_domain.ActivityLog
		if err = cursor.Decode(&activity); err != nil {
			return activity_log_domain.Response{}, err
		}
		activity.ActivityTime = activity.ActivityTime.Add(7 * time.Hour)

		wg.Add(1)
		go func(activity activity_log_domain.ActivityLog) {
			defer wg.Done()
			var admin admin_domain.Admin
			filterUser := bson.M{"_id": activity.UserID}
			_ = collectionAdmin.FindOne(ctx, filterUser).Decode(&admin)
			activity.UserID = admin.Id

			// Thêm activity vào slice activities
			activities = append(activities, activity)
		}(activity)
	}
	wg.Wait()

	var statistics activity_log_domain.Statistics
	go func() {
		statistics, _ = a.Statistics(ctx)
	}()

	activityRes := activity_log_domain.Response{
		Page:        totalPages,
		PageCurrent: int64(pageNumber),
		Statistics:  statistics,
		ActivityLog: activities,
	}

	select {
	case err = <-errCh:
		return activity_log_domain.Response{}, err
	default:
		return activityRes, nil
	}
}

func (a *activityRepository) Statistics(ctx context.Context) (activity_log_domain.Statistics, error) {
	statisticsCh := make(chan activity_log_domain.Statistics, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := statisticsCache.Get("statistics")
		if found {
			statisticsCh <- data
		}
	}()

	go func() {
		defer close(statisticsCh)
		wg.Wait()
	}()

	statisticsData := <-statisticsCh
	if !internal.IsZeroValue(statisticsData) {
		return statisticsData, nil
	}

	collection := a.database.Collection(a.collectionActivity)

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return activity_log_domain.Statistics{}, err
	}

	statistics := activity_log_domain.Statistics{
		Total: count,
	}

	statisticsCache.Set("statistics", statistics, cacheTTL)
	return statistics, nil
}
