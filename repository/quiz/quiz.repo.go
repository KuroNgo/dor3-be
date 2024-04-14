package quiz_repository

import (
	quiz_domain "clean-architecture/domain/quiz"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type quizRepository struct {
	database             *mongo.Database
	collectionLesson     string
	collectionUnit       string
	collectionVocabulary string
	collectionQuiz       string
}

func NewQuizRepository(db *mongo.Database, collectionQuiz string, collectionLesson string, collectionUnit string, collectionVocabulary string) quiz_domain.IQuizRepository {
	return &quizRepository{
		database:             db,
		collectionQuiz:       collectionQuiz,
		collectionLesson:     collectionLesson,
		collectionUnit:       collectionUnit,
		collectionVocabulary: collectionVocabulary,
	}
}
func (q *quizRepository) FetchManyByLessonID(ctx context.Context, unitID string) (quiz_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (q *quizRepository) UpdateCompleted(ctx context.Context, quizID string, isComplete int) error {
	collection := q.database.Collection(q.collectionUnit)
	objID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "is_complete", Value: isComplete},
	}}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (q *quizRepository) FetchManyByUnitID(ctx context.Context, unitID string) (quiz_domain.Response, error) {
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	idLesson2, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return quiz_domain.Response{}, err
	}

	filter := bson.M{"lesson_id": idLesson2}
	cursor, err := collectionQuiz.Find(ctx, filter)
	if err != nil {
		return quiz_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var quizs []quiz_domain.Quiz
	for cursor.Next(ctx) {
		var quiz quiz_domain.Quiz
		if err := cursor.Decode(&quiz); err != nil {
			return quiz_domain.Response{}, err
		}

		// Gắn LessonID vào đơn vị
		quiz.LessonID = idLesson2

		quizs = append(quizs, quiz)
	}

	response := quiz_domain.Response{
		Quiz: quizs,
	}
	return response, nil
}

func (q *quizRepository) FetchTenQuizButEnoughAllSkill(ctx context.Context) ([]quiz_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (q *quizRepository) FetchMany(ctx context.Context) (quiz_domain.Response, error) {
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	cursor, err := collectionQuiz.Find(ctx, bson.D{})
	if err != nil {
		return quiz_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var quiz []quiz_domain.Quiz

	for cursor.Next(ctx) {
		var q quiz_domain.Quiz
		if err := cursor.Decode(&q); err != nil {
			return quiz_domain.Response{}, err
		}
		quiz = append(quiz, q)
	}
	quizRes := quiz_domain.Response{
		Quiz: quiz,
	}
	return quizRes, nil
}

func (q *quizRepository) UpdateOne(ctx context.Context, quizID string, quiz quiz_domain.Quiz) error {
	collectionQuiz := q.database.Collection(q.collectionQuiz)
	doc, err := internal.ToDoc(quiz)
	objID, err := primitive.ObjectIDFromHex(quizID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collectionQuiz.UpdateOne(ctx, filter, update)
	return err
}

func (q *quizRepository) CreateOne(ctx context.Context, quiz *quiz_domain.Quiz) error {
	collectionQuiz := q.database.Collection(q.collectionQuiz)

	filter := bson.M{"question": quiz.Question}
	// check exists with Count Documents
	count, err := collectionQuiz.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the question did exist")
	}

	_, err = collectionQuiz.InsertOne(ctx, quiz)
	return err
}

func (q *quizRepository) DeleteOne(ctx context.Context, quizID string) error {
	collectionQuiz := q.database.Collection(q.collectionQuiz)
	objID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}
	count, err := collectionQuiz.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`the quiz is removed`)
	}
	_, err = collectionQuiz.DeleteOne(ctx, filter)
	return err
}
