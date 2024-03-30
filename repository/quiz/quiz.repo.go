package quiz_repository

import (
	quiz_domain "clean-architecture/domain/quiz"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type quizRepository struct {
	database   mongo.Database
	collection string
}

func NewQuizRepository(db mongo.Database, collection string) quiz_domain.IQuizRepository {
	return &quizRepository{
		database:   db,
		collection: collection,
	}
}

func (q *quizRepository) FetchTenQuizButEnoughAllSkill(ctx context.Context) ([]quiz_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (q *quizRepository) FetchMany(ctx context.Context) (quiz_domain.Response, error) {
	collection := q.database.Collection(q.collection)

	cursor, err := collection.Find(ctx, bson.D{})
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
	collection := q.database.Collection(q.collection)
	doc, err := internal.ToDoc(quiz)
	objID, err := primitive.ObjectIDFromHex(quizID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (q *quizRepository) CreateOne(ctx context.Context, quiz *quiz_domain.Quiz) error {
	collection := q.database.Collection(q.collection)

	filter := bson.M{"question": quiz.Question}
	// check exists with Count Documents
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the question did exist")
	}

	_, err = collection.InsertOne(ctx, quiz)
	return err
}

func (q *quizRepository) DeleteOne(ctx context.Context, quizID string) error {
	collection := q.database.Collection(q.collection)
	objID, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`the quiz is removed`)
	}
	_, err = collection.DeleteOne(ctx, filter)
	return err
}

func (q *quizRepository) UpsertOne(c context.Context, id string, quiz *quiz_domain.Quiz) (quiz_domain.Response, error) {
	collection := q.database.Collection(q.collection)
	doc, err := internal.ToDoc(quiz)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return quiz_domain.Response{}, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collection.FindOneAndUpdate(c, query, update, opts)

	var updatedPost quiz_domain.Response

	if err := res.Decode(&updatedPost); err != nil {
		return quiz_domain.Response{}, errors.New("no post with that Id exists")
	}

	return updatedPost, nil
}
