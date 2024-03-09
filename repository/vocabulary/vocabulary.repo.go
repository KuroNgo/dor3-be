package vocabulary_repository

import (
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type vocabularyRepository struct {
	database             mongo.Database
	collectionVocabulary string
	collectionLesson     string
}

func NewVocabularyRepository(db mongo.Database, collectionVocabulary string, collectionLesson string) vocabulary_domain.IVocabularyRepository {
	return &vocabularyRepository{
		database:             db,
		collectionVocabulary: collectionVocabulary,
		collectionLesson:     collectionLesson,
	}
}

func (v *vocabularyRepository) FetchByID(ctx context.Context, vocabularyID string) (*vocabulary_domain.Vocabulary, error) {
	collection := v.database.Collection(v.collectionVocabulary)

	var vocabulary vocabulary_domain.Vocabulary

	idHex, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return &vocabulary, err
	}

	err = collection.
		FindOne(ctx, bson.M{"_id": idHex}).
		Decode(&vocabulary)
	return &vocabulary, err
}

func (v *vocabularyRepository) FetchByWord(ctx context.Context, word string) ([]vocabulary_domain.Vocabulary, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	filter := bson.D{{Key: "word", Value: word}}
	// check exists with CountDocuments
	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("the word do not exist")
	}
	var vocabulary []vocabulary_domain.Vocabulary

	// Tìm kiếm tài liệu với điều kiện name
	err = collectionVocabulary.FindOne(context.Background(), bson.M{"word": word}).Decode(&vocabulary)
	if err != nil {
		return nil, err
	}

	return vocabulary, nil
}

func (v *vocabularyRepository) FetchByLesson(ctx context.Context, lessonName string) ([]vocabulary_domain.Vocabulary, error) {
	collectionLesson := v.database.Collection(v.collectionLesson)

	filter := bson.D{{Key: "name", Value: lessonName}}
	// check exists with CountDocuments
	count, err := collectionLesson.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("the lesson name do not exist")
	}
	var vocabulary []vocabulary_domain.Vocabulary

	// Tìm kiếm tài liệu với điều kiện name
	err = collectionLesson.FindOne(context.Background(), bson.M{"name": lessonName}).Decode(&vocabulary)
	if err != nil {
		return nil, err
	}

	return vocabulary, nil
}

func (v *vocabularyRepository) FetchMany(ctx context.Context) ([]vocabulary_domain.Vocabulary, error) {
	collection := v.database.Collection(v.collectionVocabulary)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var vocabulary []vocabulary_domain.Vocabulary
	err = cursor.All(ctx, &vocabulary)
	if vocabulary == nil {
		return []vocabulary_domain.Vocabulary{}, err
	}
	return vocabulary, err
}

func (v *vocabularyRepository) FetchToDeleteMany(ctx context.Context) (*[]vocabulary_domain.Vocabulary, error) {
	collection := v.database.Collection(v.collectionVocabulary)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var vocabulary *[]vocabulary_domain.Vocabulary

	err = cursor.All(ctx, vocabulary)
	if vocabulary == nil {
		return &[]vocabulary_domain.Vocabulary{}, err
	}
	return vocabulary, err
}

func (v *vocabularyRepository) UpdateOne(ctx context.Context, vocabularyID string, vocabulary vocabulary_domain.Vocabulary) error {
	collection := v.database.Collection(v.collectionVocabulary)
	doc, err := internal.ToDoc(vocabulary)
	objID, err := primitive.ObjectIDFromHex(vocabularyID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (v *vocabularyRepository) CreateOne(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionLesson := v.database.Collection(v.collectionLesson)

	filter := bson.M{"name": vocabulary.Word}
	// check exists with CountDocuments
	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the word did exist")
	}

	filterReference := bson.M{"_id": vocabulary.LessonID}
	count, err = collectionLesson.CountDocuments(ctx, filterReference)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("the lesson ID do not exist")
	}

	_, err = collectionLesson.InsertOne(ctx, vocabulary)
	return nil
}

func (v *vocabularyRepository) UpsertOne(ctx context.Context, id string, vocabulary *vocabulary_domain.Vocabulary) (*vocabulary_domain.Vocabulary, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionLesson := v.database.Collection(v.collectionLesson)

	filterReference := bson.M{"_id": vocabulary.LessonID}
	count, err := collectionLesson.CountDocuments(ctx, filterReference)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New("the lesson ID do not exist")
	}

	doc, err := internal.ToDoc(vocabulary)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collectionVocabulary.FindOneAndUpdate(ctx, query, update, opts)

	var updatePost *vocabulary_domain.Vocabulary
	if err := res.Decode(&updatePost); err != nil {
		return nil, err
	}

	return updatePost, nil
}

func (v *vocabularyRepository) DeleteOne(ctx context.Context, vocabularyID string) error {
	collection := v.database.Collection(v.collectionVocabulary)
	objID, err := primitive.ObjectIDFromHex(vocabularyID)
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
		return errors.New(`the lesson is removed`)
	}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}
