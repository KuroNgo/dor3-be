package quiz_options_repository

import (
	quiz_options_domain "clean-architecture/domain/quiz_options"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type quizOptionsRepository struct {
	database           *mongo.Database
	collectionQuestion string
	collectionOptions  string
}

func (q *quizOptionsRepository) FetchManyByQuestionID(ctx context.Context, questionID string) (quiz_options_domain.Response, error) {
	collectionOptions := q.database.Collection(q.collectionOptions)
	idQuestion, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return quiz_options_domain.Response{}, err
	}

	filter := bson.M{"question_id": idQuestion}
	cursor, err := collectionOptions.Find(ctx, filter)
	if err != nil {
		return quiz_options_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var options []quiz_options_domain.QuizOptions
	for cursor.Next(ctx) {
		var option quiz_options_domain.QuizOptions
		if err = cursor.Decode(&option); err != nil {
			return quiz_options_domain.Response{}, err
		}

		// Gắn CourseID vào bài học
		option.QuestionID = idQuestion
		options = append(options, option)
	}

	response := quiz_options_domain.Response{
		QuizOptions: options,
	}

	return response, nil
}

func (q *quizOptionsRepository) UpdateOne(ctx context.Context, quizOptions *quiz_options_domain.QuizOptions) (*mongo.UpdateResult, error) {
	collection := q.database.Collection(q.collectionQuestion)

	filter := bson.D{{Key: "_id", Value: quizOptions.ID}}
	update := bson.M{
		"$set": bson.M{
			"question_id": quizOptions.QuestionID,
			"content":     quizOptions.Content,
			"update_at":   quizOptions.UpdateAt,
			"who_update":  quizOptions.WhoUpdate,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (q *quizOptionsRepository) CreateOne(ctx context.Context, quizOptions *quiz_options_domain.QuizOptions) error {
	collectionOptions := q.database.Collection(q.collectionOptions)
	collectionQuestion := q.database.Collection(q.collectionQuestion)

	filterQuestionID := bson.M{"question_id": quizOptions.QuestionID}
	countLessonID, err := collectionQuestion.CountDocuments(ctx, filterQuestionID)
	if err != nil {
		return err
	}

	if countLessonID == 0 {
		return errors.New("the question ID do not exist")
	}

	_, err = collectionOptions.InsertOne(ctx, quizOptions)
	return nil
}

func (q *quizOptionsRepository) DeleteOne(ctx context.Context, optionsID string) error {
	collectionOptions := q.database.Collection(q.collectionOptions)
	objID, err := primitive.ObjectIDFromHex(optionsID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionOptions.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`quiz answer is removed`)
	}

	_, err = collectionOptions.DeleteOne(ctx, filter)
	return err
}

func NewQuizOptionsRepository(db *mongo.Database, collectionQuestion string, collectionOptions string) quiz_options_domain.IQuizOptionRepository {
	return &quizOptionsRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionOptions:  collectionOptions,
	}
}
