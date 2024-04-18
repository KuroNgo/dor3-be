package exam_question

import (
	exam_question_domain "clean-architecture/domain/exam_question"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type examQuestionRepository struct {
	database           mongo.Database
	collectionQuestion string
	collectionExam     string
}

func NewExamQuestionRepository(db mongo.Database, collectionQuestion string, collectionExam string) exam_question_domain.IExamQuestionRepository {
	return &examQuestionRepository{
		database:           db,
		collectionQuestion: collectionQuestion,
		collectionExam:     collectionExam,
	}
}

func (e *examQuestionRepository) FetchMany(ctx context.Context, page string) (exam_question_domain.Response, error) {
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return exam_question_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	cursor, err := collectionQuestion.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	var questions []exam_question_domain.ExamQuestion
	for cursor.Next(ctx) {
		var question exam_question_domain.ExamQuestion
		if err = cursor.Decode(&question); err != nil {
			return exam_question_domain.Response{}, err
		}

		questions = append(questions, question)
	}
	questionsRes := exam_question_domain.Response{
		ExamQuestion: questions,
	}
	return questionsRes, nil
}

func (e *examQuestionRepository) FetchManyByExamID(ctx context.Context, examID string) (exam_question_domain.Response, error) {
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	idExam, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	filter := bson.M{"exam_id": idExam}
	cursor, err := collectionQuestion.Find(ctx, filter)
	if err != nil {
		return exam_question_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var questions []exam_question_domain.ExamQuestion
	for cursor.Next(ctx) {
		var question exam_question_domain.ExamQuestion
		if err = cursor.Decode(&question); err != nil {
			return exam_question_domain.Response{}, err
		}

		question.ExamID = idExam
		questions = append(questions, question)
	}

	questionsRes := exam_question_domain.Response{
		ExamQuestion: questions,
	}

	return questionsRes, nil
}

func (e *examQuestionRepository) UpdateOne(ctx context.Context, examQuestion *exam_question_domain.ExamQuestion) (*mongo.UpdateResult, error) {
	collection := e.database.Collection(e.collectionQuestion)

	filter := bson.D{{Key: "_id", Value: examQuestion.ID}}
	update := bson.M{
		"$set": bson.M{
			"exam_id":        examQuestion.ExamID,
			"content":        examQuestion.Content,
			"level":          examQuestion.Level,
			"filename":       examQuestion.Filename,
			"audio_duration": examQuestion.AudioDuration,
			"update_at":      examQuestion.UpdateAt,
			"who_update":     examQuestion.WhoUpdate,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *examQuestionRepository) CreateOne(ctx context.Context, examQuestion *exam_question_domain.ExamQuestion) error {
	collectionQuestion := e.database.Collection(e.collectionQuestion)
	collectionExam := e.database.Collection(e.collectionExam)

	filterExamID := bson.M{"exam_id": examQuestion.ExamID}
	countLessonID, err := collectionExam.CountDocuments(ctx, filterExamID)
	if err != nil {
		return err
	}

	if countLessonID == 0 {
		return errors.New("the examID do not exist")
	}

	_, err = collectionQuestion.InsertOne(ctx, examQuestion)
	return nil
}

func (e *examQuestionRepository) DeleteOne(ctx context.Context, examID string) error {
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	objID, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionQuestion.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exam answer is removed`)
	}

	_, err = collectionQuestion.DeleteOne(ctx, filter)
	return err
}
