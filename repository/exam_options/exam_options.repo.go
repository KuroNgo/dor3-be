package exam_options_repository

import (
	exam_options_domain "clean-architecture/domain/exam_options"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type examOptionsRepository struct {
	database           *mongo.Database
	collectionQuestion string
	collectionOptions  string
}

func NewExamOptionsRepository(db *mongo.Database, collectionQuestion string, collectionOptions string) exam_options_domain.IExamOptionRepository {
	return &examOptionsRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionOptions:  collectionOptions,
	}
}

func (e *examOptionsRepository) FetchManyByQuestionID(ctx context.Context, questionID string) (exam_options_domain.Response, error) {
	collectionOptions := e.database.Collection(e.collectionOptions)
	idQuestion, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return exam_options_domain.Response{}, err
	}

	filter := bson.M{"question_id": idQuestion}
	cursor, err := collectionOptions.Find(ctx, filter)
	if err != nil {
		return exam_options_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var options []exam_options_domain.ExamOptions

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var option exam_options_domain.ExamOptions
			if err = cursor.Decode(&option); err != nil {
				return
			}

			// Gắn CourseID vào bài học
			option.QuestionID = idQuestion
			options = append(options, option)
		}
	}()

	wg.Done()

	response := exam_options_domain.Response{
		ExamOptions: options,
	}

	return response, nil
}

func (e *examOptionsRepository) UpdateOne(ctx context.Context, examOptions *exam_options_domain.ExamOptions) (*mongo.UpdateResult, error) {
	collection := e.database.Collection(e.collectionQuestion)

	filter := bson.D{{Key: "_id", Value: examOptions.ID}}
	update := bson.M{
		"$set": bson.M{
			"question_id": examOptions.QuestionID,
			"content":     examOptions.Content,
			"update_at":   examOptions.UpdateAt,
			"who_update":  examOptions.WhoUpdate,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *examOptionsRepository) CreateOne(ctx context.Context, examOptions *exam_options_domain.ExamOptions) error {
	collectionOptions := e.database.Collection(e.collectionOptions)
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	filterQuestionID := bson.M{"question_id": examOptions.QuestionID}
	countLessonID, err := collectionQuestion.CountDocuments(ctx, filterQuestionID)
	if err != nil {
		return err
	}

	if countLessonID == 0 {
		return errors.New("the question ID do not exist")
	}

	_, err = collectionOptions.InsertOne(ctx, examOptions)
	return nil
}

func (e *examOptionsRepository) DeleteOne(ctx context.Context, optionsID string) error {
	collectionOptions := e.database.Collection(e.collectionOptions)
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
		return errors.New(`exam answer is removed`)
	}

	_, err = collectionOptions.DeleteOne(ctx, filter)
	return err
}
