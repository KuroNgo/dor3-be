package quiz

import (
	quiz_domain "clean-architecture/domain/quiz"
	"clean-architecture/infrastructor/mongo"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (q *quizRepository) Fetch(ctx context.Context) ([]quiz_domain.Quiz, error) {
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

func (q *quizRepository) FetchToDelete(ctx context.Context) (*[]quiz_domain.Quiz, error) {
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

func (q *quizRepository) Update(ctx context.Context, quizID string, quiz quiz_domain.Quiz) error {
	collection := q.database.Collection(q.collection)
	filter := bson.M{"_id": quizID}

	update := bson.M{
		"$set": bson.M{
			"Question":      quiz.Question,
			"Options":       quiz.Options,
			"CorrectAnswer": quiz.CorrectAnswer,
			"QuestionType":  quiz.QuestionType,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (q *quizRepository) Create(ctx context.Context, quiz *quiz_domain.Quiz) error {
	collection := q.database.Collection(q.collection)
	_, err := collection.InsertOne(ctx, quiz)
	return err
}

func (q *quizRepository) Delete(ctx context.Context, quizID string) error {
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
