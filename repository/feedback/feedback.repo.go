package feedback_repository

import (
	feedback_domain "clean-architecture/domain/feedback"
	user_domain "clean-architecture/domain/user"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"sync"
	"time"
)

type feedbackRepository struct {
	database           *mongo.Database
	collectionFeedback string
	collectionUser     string
}

func NewFeedbackRepository(db *mongo.Database, collectionFeedback string, collectionUser string) feedback_domain.IFeedbackRepository {
	return &feedbackRepository{
		database:           db,
		collectionFeedback: collectionFeedback,
		collectionUser:     collectionUser,
	}
}

var (
	wg sync.WaitGroup
)

func (f *feedbackRepository) FetchManyInAdmin(ctx context.Context, page string) (feedback_domain.Response, error) {
	collection := f.database.Collection(f.collectionFeedback)
	collectionUser := f.database.Collection(f.collectionUser)

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

	var feedbacks []feedback_domain.FeedbackResponse
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var feedback feedback_domain.Feedback
			if err = cursor.Decode(&feedback); err != nil {
				return
			}

			var user user_domain.User
			filterUser := bson.M{"_id": feedback.UserID}
			_ = collectionUser.FindOne(ctx, filterUser).Decode(&user)

			feedbackRes := feedback_domain.FeedbackResponse{
				ID:            feedback.ID,
				User:          user,
				Feeling:       feedback.Feeling,
				Content:       feedback.Content,
				Title:         feedback.Title,
				IsSeen:        feedback.IsSeen,
				SeenAt:        feedback.SeenAt,
				IsLoveWeb:     feedback.IsLoveWeb,
				SubmittedDate: feedback.SubmittedDate,
			}

			feedbacks = append(feedbacks, feedbackRes)
		}
	}()

	wg.Wait()

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

func (f *feedbackRepository) FetchByUserIDInAdmin(ctx context.Context, userID string, page string) (feedback_domain.Response, error) {
	collection := f.database.Collection(f.collectionFeedback)
	collectionUser := f.database.Collection(f.collectionUser)

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

	var feedbacks []feedback_domain.FeedbackResponse
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var feedback feedback_domain.Feedback
			if err := cursor.Decode(&feedback); err != nil {
				return
			}

			feedback.UserID = idUser
			var user user_domain.User
			filterUser := bson.M{"_id": feedback.UserID}
			_ = collectionUser.FindOne(ctx, filterUser).Decode(&user)

			var feedbackRes feedback_domain.FeedbackResponse
			feedbackRes.ID = feedback.ID
			feedbackRes.User = user
			feedbackRes.Feeling = feedback.Feeling
			feedbackRes.Content = feedback.Content
			feedbackRes.Title = feedback.Title
			feedbackRes.IsSeen = feedback.IsSeen
			feedbackRes.SeenAt = feedback.SeenAt
			feedbackRes.IsLoveWeb = feedback.IsLoveWeb
			feedbackRes.SubmittedDate = feedback.SubmittedDate

			feedbacks = append(feedbacks, feedbackRes)
		}
	}()
	wg.Wait()

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
		Feedback:    feedbacks,
		Statistics:  statistics,
	}

	return feedbackRes, nil
}

func (f *feedbackRepository) FetchBySubmittedDateInAdmin(ctx context.Context, date string, page string) (feedback_domain.Response, error) {
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
	}

	return feedbackRes, nil
}

func (f *feedbackRepository) CreateOneInUser(ctx context.Context, feedback *feedback_domain.Feedback) error {
	collectionFeedback := f.database.Collection(f.collectionFeedback)
	_, err := collectionFeedback.InsertOne(ctx, feedback)
	return err
}

func (f *feedbackRepository) UpdateSeenInAdmin(ctx context.Context, id string, isSeen int) error {
	collection := f.database.Collection(f.collectionFeedback)

	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: ID}}
	update := bson.M{
		"$set": bson.M{
			"is_seen": isSeen,
			"seen_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return err
	}
	return nil
}

func (f *feedbackRepository) DeleteOneInAdmin(ctx context.Context, feedbackID string) error {
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
	// Lấy collection từ database
	collection := f.database.Collection(f.collectionFeedback)

	// Đếm tổng số tài liệu trong collection
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return feedback_domain.Statistics{}, err
	}

	// Tìm tất cả tài liệu trong collection
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return feedback_domain.Statistics{}, err
	}
	// Đảm bảo con trỏ được đóng sau khi sử dụng
	defer cursor.Close(ctx)

	// Khởi tạo các biến đếm
	var (
		countIsLove       int32 = 0
		countIsSeen       int32 = 0
		countIsNotSeen    int32 = 0
		countSad          int32 = 0
		countHappy        int32 = 0
		countDisappointed int32 = 0
		countGood         int32 = 0
	)

	// Lặp qua các tài liệu và cập nhật các biến đếm
	for cursor.Next(ctx) {
		var feedback feedback_domain.Feedback
		if err = cursor.Decode(&feedback); err != nil {
			return feedback_domain.Statistics{}, err
		}

		// Kiểm tra và cập nhật các biến đếm tương ứng
		if feedback.IsLoveWeb == 1 {
			countIsLove++
		}
		if feedback.IsSeen == 1 {
			countIsSeen++
		} else {
			countIsNotSeen++
		}

		// Sử dụng switch-case để kiểm tra cảm xúc và cập nhật biến đếm tương ứng
		switch feedback.Feeling {
		case "sad":
			countSad++
		case "happy":
			countHappy++
		case "disappointed":
			countDisappointed++
		case "good":
			countGood++
		}
	}

	// Kiểm tra xem có lỗi nào xảy ra trong quá trình lặp qua con trỏ không
	if err := cursor.Err(); err != nil {
		return feedback_domain.Statistics{}, err
	}

	// Tính toán phần trăm các cảm xúc
	var percentSad, percentHappy, percentDisappointed, percentGood float32
	if count > 0 {
		percentSad = float32(countSad) / float32(count) * 100
		percentHappy = float32(countHappy) / float32(count) * 100
		percentDisappointed = float32(countDisappointed) / float32(count) * 100
		percentGood = float32(countGood) / float32(count) * 100
	}

	// Tạo đối tượng feedbackRes để trả về
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

	// Trả về kết quả
	return feedbackRes, nil
}
