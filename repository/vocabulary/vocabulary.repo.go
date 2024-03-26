package vocabulary_repository

import (
	unit_domain "clean-architecture/domain/unit"
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

func (v *vocabularyRepository) FetchByIdUnit(ctx context.Context, idUnit string) (vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	idUnit2, err := primitive.ObjectIDFromHex(idUnit)
	filter := bson.M{"unit_id": idUnit2}

	cursor, err := collectionVocabulary.Find(ctx, filter)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var vocabularies []vocabulary_domain.Vocabulary
	// Lặp qua các kết quả và giải mã vào slice units
	for cursor.Next(ctx) {
		var vocabulary vocabulary_domain.Vocabulary
		if err = cursor.Decode(&vocabulary); err != nil {
			return vocabulary_domain.Response{}, err
		}

		var unit unit_domain.Unit
		err = collectionUnit.FindOne(ctx, bson.M{"_id": idUnit2}).Decode(&unit)
		if err != nil {
			return vocabulary_domain.Response{}, err
		}

		vocabulary.UnitID = idUnit2

		vocabularies = append(vocabularies, vocabulary)
	}

	// Tạo và trả về phản hồi với dữ liệu units và số lượng tài liệu trong collection bài học
	response := vocabulary_domain.Response{
		Vocabulary: vocabularies,
	}
	return response, nil
}

func (v *vocabularyRepository) FetchByWord(ctx context.Context, word string) (vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	filter := bson.M{"word": primitive.Regex{Pattern: word, Options: "i"}}
	var vocabularies []vocabulary_domain.Vocabulary

	// Tìm kiếm tài liệu với điều kiện name
	cursor, err := collectionVocabulary.Find(ctx, filter)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}
	for cursor.Next(ctx) {
		var vocabulary vocabulary_domain.Vocabulary
		if err = cursor.Decode(&vocabulary); err != nil {
			return vocabulary_domain.Response{}, err
		}

		// Thêm lesson vào slice lessons
		vocabularies = append(vocabularies, vocabulary)
	}
	err = cursor.All(ctx, &vocabularies)
	vocabularyRes := vocabulary_domain.Response{
		Vocabulary: vocabularies,
	}

	return vocabularyRes, nil
}

func (v *vocabularyRepository) FetchByLesson(ctx context.Context, unitName string) (vocabulary_domain.Response, error) {
	collectionUnit := v.database.Collection(v.collectionUnit)
	var vocabulary vocabulary_domain.Response

	// Tìm kiếm tài liệu với điều kiện name
	cursor, err := collectionUnit.Find(ctx, bson.M{"name": unitName})
	if err != nil {
		return vocabulary_domain.Response{}, err
	}
	err = cursor.All(ctx, &vocabulary)

	return vocabulary, nil
}

func (v *vocabularyRepository) FetchMany(ctx context.Context) (vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	cursor, err := collectionVocabulary.Find(ctx, bson.D{})
	if err != nil {
		return vocabulary_domain.Response{}, err
	}
	defer cursor.Close(ctx)
	var vocabularies []vocabulary_domain.Vocabulary

	for cursor.Next(ctx) {
		var vocabulary vocabulary_domain.Vocabulary
		if err = cursor.Decode(&vocabulary); err != nil {
			return vocabulary_domain.Response{}, err
		}

		var unit unit_domain.Unit
		err = collectionUnit.FindOne(ctx, bson.M{"_id": vocabulary.UnitID}).Decode(&unit)
		if err != nil {
			return vocabulary_domain.Response{}, err
		}

		// Gắn tên của course vào lesson
		vocabulary.UnitID = unit.ID

		// Thêm lesson vào slice lessons
		vocabularies = append(vocabularies, vocabulary)
	}

	var vocabulary vocabulary_domain.Response
	err = cursor.All(ctx, &vocabulary)
	vocabularyRes := vocabulary_domain.Response{
		Vocabulary: vocabularies,
	}
	return vocabularyRes, err
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

	filter := bson.M{"word": vocabulary.Word, "unit_id": vocabulary.UnitID}
	filterReference := bson.M{"_id": vocabulary.UnitID}

	count, err := collectionUnit.CountDocuments(ctx, filterReference)
	if err != nil {
		return err
	}

	// check exists with CountDocuments
	count, err = collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the word in unit did exist")
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
