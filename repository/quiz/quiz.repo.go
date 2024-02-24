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

func (q *quizRepository) FetchByID(ctx context.Context, quizID string) (*quiz_domain.Quiz, error) {
	collection := q.database.Collection(q.collection)

	var quiz quiz_domain.Quiz

	idHex, err := primitive.ObjectIDFromHex(quizID)
	if err != nil {
		return &quiz, err
	}

	err = collection.
		FindOne(ctx, bson.M{"_id": idHex}).
		Decode(&quiz)
	return &quiz, err
}

func (q *quizRepository) FetchMany(ctx context.Context) ([]quiz_domain.Quiz, error) {
	collection := q.database.Collection(q.collection)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var quiz []quiz_domain.Quiz

	err = cursor.All(ctx, &quiz)
	if quiz == nil {
		return []quiz_domain.Quiz{}, err
	}

	return quiz, err
}

func (q *quizRepository) FetchToDeleteMany(ctx context.Context) (*[]quiz_domain.Quiz, error) {
	collection := q.database.Collection(q.collection)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var quiz []quiz_domain.Quiz

	err = cursor.All(ctx, &quiz)
	if quiz == nil {
		return &[]quiz_domain.Quiz{}, err
	}

	return &quiz, err
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

func (q *quizRepository) CreateOne(ctx context.Context, quiz *quiz_domain.Input) error {
	collection := q.database.Collection(q.collection)
	_, err := collection.InsertOne(ctx, quiz)
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
	_, err = collection.DeleteOne(ctx, filter)
	return err
}

func (q *quizRepository) UpsertOne(c context.Context, id string, quiz *quiz_domain.Quiz) (*quiz_domain.Response, error) {
	collection := q.database.Collection(q.collection)
	doc, err := internal.ToDoc(quiz)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collection.FindOneAndUpdate(c, query, update, opts)

	var updatedPost *quiz_domain.Response

	if err := res.Decode(&updatedPost); err != nil {
		return nil, errors.New("no post with that Id exists")
	}

	return updatedPost, nil
}
