package feedback_repository

import (
	feedback_domain "clean-architecture/domain/feedback"
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
	perPage := 1
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
	for cursor.Next(ctx) {
		var feedback feedback_domain.Feedback
		if err := cursor.Decode(&feedback); err != nil {
			return feedback_domain.Response{}, err
		}

		feedbacks = append(feedbacks, feedback)
	}

	feedbackRes := feedback_domain.Response{
		Page:     cal,
		Total:    count,
		Feedback: feedbacks,
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
		Total:    count,
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
		Page:     cal,
		Total:    count,
		Feedback: feedbacks,
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
