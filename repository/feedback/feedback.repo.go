package feedback_repository

import (
	feedback_domain "clean-architecture/domain/feedback"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type feedbackRepository struct {
	database           *mongo.Database
	collectionFeedback string
}

func NewFeedbackRepository(db *mongo.Database, collectionFeedback string) feedback_domain.IFeedbackRepository {
	return &feedbackRepository{
		database:           db,
		collectionFeedback: collectionFeedback,
	}
}

func (f *feedbackRepository) FetchMany(ctx context.Context, page string) (feedback_domain.Response, error) {
	collection := f.database.Collection(f.collectionFeedback)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return feedback_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 10
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return feedback_domain.Response{}, err
	}

	cal1 := count / int64(perPage)
	cal2 := count % int64(perPage)
	var cal int64 = 0
	if cal2 != 0 {
		cal = cal1 + 1
	}

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return feedback_domain.Response{}, err
	}

	var feedbacks []feedback_domain.Feedback
	internal.Wg.Add(1)

	go func() {
		defer internal.Wg.Done()
		for cursor.Next(ctx) {
			var feedback feedback_domain.Feedback
			if err := cursor.Decode(&feedback); err != nil {
				return
			}

			feedbacks = append(feedbacks, feedback)
		}
	}()

	internal.Wg.Wait()

	var statisticsCh = make(chan feedback_domain.Statistics)
	go func() {
		defer close(statisticsCh)
		statistics, _ := f.Statistics(ctx)
		statisticsCh <- statistics
	}()

	statistics := <-statisticsCh
	feedbackRes := feedback_domain.Response{
		Page:        cal,
		CurrentPage: int64(pageNumber),
		Statistics:  statistics,
		Feedback:    feedbacks,
	}

	return feedbackRes, nil
}

func (f *feedbackRepository) FetchByUserID(ctx context.Context, userID string, page string) (feedback_domain.Response, error) {
	collection := f.database.Collection(f.collectionFeedback)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return feedback_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 1
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return feedback_domain.Response{}, err
	}
	filter := bson.M{"user_id": idUser}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return feedback_domain.Response{}, err
	}

	cal1 := count / int64(perPage)
	cal2 := count % int64(perPage)
	var cal int64 = 0
	if cal2 != 0 {
		cal = cal1 + 1
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return feedback_domain.Response{}, err
	}

	var feedbacks []feedback_domain.Feedback
	for cursor.Next(ctx) {
		var feedback feedback_domain.Feedback
		if err := cursor.Decode(&feedback); err != nil {
			return feedback_domain.Response{}, err
		}

		feedback.UserID = idUser

		feedbacks = append(feedbacks, feedback)
	}

	feedbackRes := feedback_domain.Response{
		Page:     cal,
		Feedback: feedbacks,
	}

	return feedbackRes, nil
}

func (f *feedbackRepository) FetchBySubmittedDate(ctx context.Context, date string, page string) (feedback_domain.Response, error) {
	collection := f.database.Collection(f.collectionFeedback)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return feedback_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 1
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	submittedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return feedback_domain.Response{}, err
	}
	filter := bson.M{"submitted_date": submittedDate}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return feedback_domain.Response{}, err
	}

	cal1 := count / int64(perPage)
	cal2 := count % int64(perPage)
	var cal int64 = 0
	if cal2 != 0 {
		cal = cal1 + 1
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return feedback_domain.Response{}, err
	}

	var feedbacks []feedback_domain.Feedback
	for cursor.Next(ctx) {
		var feedback feedback_domain.Feedback
		if err := cursor.Decode(&feedback); err != nil {
			return feedback_domain.Response{}, err
		}

		feedback.SubmittedDate = submittedDate
		feedbacks = append(feedbacks, feedback)
	}

	feedbackRes := feedback_domain.Response{
		CurrentPage: int64(pageNumber),
		Page:        cal,
		Feedback:    feedbacks,
	}

	return feedbackRes, nil
}

func (f *feedbackRepository) CreateOneByUser(ctx context.Context, feedback *feedback_domain.Feedback) error {
	collectionFeedback := f.database.Collection(f.collectionFeedback)
	_, err := collectionFeedback.InsertOne(ctx, feedback)
	return err
}

func (f *feedbackRepository) DeleteOneByAdmin(ctx context.Context, feedbackID string) error {
	collectionFeedback := f.database.Collection(f.collectionFeedback)

	objID, err := primitive.ObjectIDFromHex(feedbackID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}

	count, err := collectionFeedback.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`the feedback is removed or have not exist`)
	}

	_, err = collectionFeedback.DeleteOne(ctx, filter)
	return err
}

func (f *feedbackRepository) Statistics(ctx context.Context) (feedback_domain.Statistics, error) {
	collection := f.database.Collection(f.collectionFeedback)

	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return feedback_domain.Statistics{}, err
	}

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return feedback_domain.Statistics{}, err
	}

	var feedbacks []feedback_domain.Feedback
	var feedback feedback_domain.Feedback

	var (
		countIsLove       int32 = 0
		countIsSeen       int32 = 0
		countIsNotSeen    int32 = 0
		countSad          int32 = 0 //sad, happy, disappointed, good
		countHappy        int32 = 0
		countDisappointed int32 = 0
		countGood         int32 = 0
	)

	for cursor.Next(ctx) {
		if err = cursor.Decode(&feedback); err != nil {
			return feedback_domain.Statistics{}, err
		}

		if feedback.IsLoveWeb == 1 {
			countIsLove++
		}
		if feedback.IsSeen == 1 {
			countIsSeen++
		}
		if feedback.IsSeen == 0 {
			countIsNotSeen++
		}
		if feedback.Feeling == "sad" {
			countSad++
		}
		if feedback.Feeling == "happy" {
			countHappy++
		}
		if feedback.Feeling == "disappointed" {
			countDisappointed++
		}
		if feedback.Feeling == "good" {
			countGood++
		}

		feedbacks = append(feedbacks, feedback)
	}

	var percentSad, percentHappy, percentDisappointed, percentGood float32
	if countSad != 0 {
		percentSad = float32((float64(countSad) / float64(count)) * 100)
	}
	if countHappy != 0 {
		percentHappy = float32((float64(countHappy) / float64(count)) * 100)
	}
	if countDisappointed != 0 {
		percentDisappointed = float32((float64(countDisappointed) / float64(count)) * 100)
	}
	if countGood != 0 {
		percentGood = float32((float64(countGood) / float64(count)) * 100)
	}

	feedbackRes := feedback_domain.Statistics{
		Total:             count,
		TotalIsLoveWeb:    countIsLove,
		TotalIsSeen:       countIsSeen,
		TotalIsNotSeen:    countIsNotSeen,
		TotalFeeling:      int32(count),
		CountSad:          percentSad,
		CountHappy:        percentHappy,
		CountDisappointed: percentDisappointed,
		CountGood:         percentGood,
	}

	return feedbackRes, nil
}
