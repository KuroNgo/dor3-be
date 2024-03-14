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
	collectionUnit       string
}

func NewVocabularyRepository(db mongo.Database, collectionVocabulary string, collectionUnit string) vocabulary_domain.IVocabularyRepository {
	return &vocabularyRepository{
		database:             db,
		collectionVocabulary: collectionVocabulary,
		collectionUnit:       collectionUnit,
	}
}

func (v *vocabularyRepository) FetchByWord(ctx context.Context, word string) ([]vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	filter := bson.M{"word": primitive.Regex{Pattern: word, Options: "i"}}
	var vocabularies []vocabulary_domain.Response

	// Tìm kiếm tài liệu với điều kiện name
	cursor, err := collectionVocabulary.Find(ctx, filter)
	err = cursor.All(ctx, &vocabularies)
	if vocabularies == nil {
		return []vocabulary_domain.Response{}, err
	}

	return vocabularies, nil
}

func (v *vocabularyRepository) FetchByLesson(ctx context.Context, unitName string) ([]vocabulary_domain.Response, error) {
	collectionUnit := v.database.Collection(v.collectionUnit)
	var vocabulary []vocabulary_domain.Response

	// Tìm kiếm tài liệu với điều kiện name
	cursor, err := collectionUnit.Find(ctx, bson.M{"name": unitName})
	err = cursor.All(ctx, &vocabulary)
	if vocabulary == nil {
		return []vocabulary_domain.Response{}, err
	}

	return vocabulary, nil
}

func (v *vocabularyRepository) FetchMany(ctx context.Context) ([]vocabulary_domain.Response, error) {
	collection := v.database.Collection(v.collectionVocabulary)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var vocabulary []vocabulary_domain.Response
	err = cursor.All(ctx, &vocabulary)
	if vocabulary == nil {
		return []vocabulary_domain.Response{}, err
	}
	return vocabulary, err
}

func (v *vocabularyRepository) FetchToDeleteMany(ctx context.Context) (*[]vocabulary_domain.Response, error) {
	collection := v.database.Collection(v.collectionVocabulary)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var vocabulary *[]vocabulary_domain.Response

	err = cursor.All(ctx, vocabulary)
	if vocabulary == nil {
		return &[]vocabulary_domain.Response{}, err
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
	collectionUnit := v.database.Collection(v.collectionUnit)

	filter := bson.M{"name": vocabulary.Word}
	// check exists with CountDocuments
	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the word did exist")
	}

	filterReference := bson.M{"_id": vocabulary.UnitID}
	count, err = collectionUnit.CountDocuments(ctx, filterReference)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("the unit ID do not exist")
	}

	_, err = collectionVocabulary.InsertOne(ctx, vocabulary)
	return nil
}

func (v *vocabularyRepository) UpsertOne(ctx context.Context, id string, vocabulary *vocabulary_domain.Vocabulary) (*vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	filterReference := bson.M{"_id": vocabulary.UnitID}
	count, err := collectionUnit.CountDocuments(ctx, filterReference)
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

	var updatePost *vocabulary_domain.Response
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
