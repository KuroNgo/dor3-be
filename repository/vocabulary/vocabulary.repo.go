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
	"strconv"
)

type vocabularyRepository struct {
	database             mongo.Database
	collectionVocabulary string
	collectionMean       string
	collectionMark       string
	collectionUnit       string
}

func NewVocabularyRepository(db mongo.Database, collectionVocabulary string, collectionMean string, collectionMark string, collectionUnit string) vocabulary_domain.IVocabularyRepository {
	return &vocabularyRepository{
		database:             db,
		collectionVocabulary: collectionVocabulary,
		collectionMean:       collectionMean,
		collectionMark:       collectionMark,
		collectionUnit:       collectionUnit,
	}
}

func (v *vocabularyRepository) FindUnitIDByUnitName(ctx context.Context, unitName string) (primitive.ObjectID, error) {
	collectionUnit := v.database.Collection(v.collectionUnit)

	filter := bson.M{"name": unitName}
	var data struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	err := collectionUnit.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return data.Id, nil
}

func (v *vocabularyRepository) GetAllVocabulary(ctx context.Context) ([]string, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	var vocabularies []string

	cursor, err := collectionVocabulary.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result bson.M
		if err = cursor.Decode(&result); err != nil {
			return nil, err
		}
		word, ok := result["word"].(string)
		if !ok {
			return nil, errors.New("failed to parse word from result")
		}
		vocabularies = append(vocabularies, word)
	}

	return vocabularies, nil
}

func (v *vocabularyRepository) CreateOneByNameUnit(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	filter := bson.M{"word": vocabulary.Word, "unit_id": vocabulary.UnitID}

	filterParent := bson.M{"_id": vocabulary.UnitID}
	countParent, err := collectionUnit.CountDocuments(ctx, filterParent)
	if err != nil {
		return err
	}
	if countParent == 0 {
		return errors.New("parent unit not found")
	}

	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the vocabulary already exists in the lesson")
	}

	_, err = collectionVocabulary.InsertOne(ctx, vocabulary)
	if err != nil {
		return err
	}
	return nil
}

func (v *vocabularyRepository) FetchByIdUnit(ctx context.Context, idUnit string) (vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	idUnit2, err := primitive.ObjectIDFromHex(idUnit)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}

	filter := bson.M{"unit_id": idUnit2}

	cursor, err := collectionVocabulary.Find(ctx, filter)
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
		vocabulary.UnitID = idUnit2
		vocabularies = append(vocabularies, vocabulary)
	}

	if len(vocabularies) == 0 {
		return vocabulary_domain.Response{}, errors.New("no vocabulary found for the provided unit_id")
	}

	response := vocabulary_domain.Response{
		Vocabulary: vocabularies,
	}
	return response, nil
}

func (v *vocabularyRepository) FetchByWord(ctx context.Context, word string) (vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	regex := primitive.Regex{Pattern: word, Options: "i"}
	filter := bson.M{"word": bson.M{"$regex": regex}}

	var limit int64 = 10

	cursor, err := collectionVocabulary.Find(ctx, filter, &options.FindOptions{Limit: &limit})
	if err != nil {
		return vocabulary_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var vocabularies []vocabulary_domain.Vocabulary
	if err := cursor.All(ctx, &vocabularies); err != nil {
		return vocabulary_domain.Response{}, err
	}

	vocabularyRes := vocabulary_domain.Response{
		Vocabulary: vocabularies,
	}

	return vocabularyRes, nil
}

func (v *vocabularyRepository) FetchByWord2(ctx context.Context, word string) (vocabulary_domain.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (v *vocabularyRepository) FetchByLesson(ctx context.Context, lessonName string) (vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	regex := primitive.Regex{Pattern: lessonName, Options: "i"}
	filter := bson.M{"field_of_it": bson.M{"$regex": regex}}

	var limit int64 = 10

	cursor, err := collectionVocabulary.Find(ctx, filter, &options.FindOptions{Limit: &limit})
	if err != nil {
		return vocabulary_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var vocabularies []vocabulary_domain.Vocabulary
	if err := cursor.All(ctx, &vocabularies); err != nil {
		return vocabulary_domain.Response{}, err
	}

	vocabularyRes := vocabulary_domain.Response{
		Vocabulary: vocabularies,
	}

	return vocabularyRes, nil
}

func (v *vocabularyRepository) FetchMany(ctx context.Context, page string) (vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return vocabulary_domain.Response{}, errors.New("invalid page number")
	}

	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	cursor, err := collectionVocabulary.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var vocabularies []vocabulary_domain.Vocabulary

	for cursor.Next(ctx) {
		var vocabulary vocabulary_domain.Vocabulary

		if err := cursor.Decode(&vocabulary); err != nil {
			return vocabulary_domain.Response{}, err
		}

		var unit unit_domain.Unit
		if err := collectionUnit.FindOne(ctx, bson.M{"_id": vocabulary.UnitID}).Decode(&unit); err != nil {
			return vocabulary_domain.Response{}, err
		}

		vocabulary.UnitID = unit.ID
		vocabularies = append(vocabularies, vocabulary)
	}

	vocabularyRes := vocabulary_domain.Response{
		Vocabulary: vocabularies,
	}

	return vocabularyRes, nil
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

func (v *vocabularyRepository) UpdateOneAudio(c context.Context, vocabularyID string, linkURL string) error {
	collection := v.database.Collection(v.collectionVocabulary)
	objID, err := primitive.ObjectIDFromHex(vocabularyID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.M{
		"$set": bson.M{
			"linkURL": linkURL,
		},
	}

	_, err = collection.UpdateOne(c, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyRepository) CreateOne(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	filter := bson.M{"word": vocabulary.Word, "unit_id": vocabulary.UnitID}
	filterReference := bson.M{"_id": vocabulary.UnitID}

	countParent, err := collectionUnit.CountDocuments(ctx, filterReference)
	if err != nil {
		return err
	}

	// check exists with CountDocuments
	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the word in unit did exist")
	}
	if countParent == 0 {
		return errors.New("the unit ID do not exist")
	}

	_, err = collectionVocabulary.InsertOne(ctx, vocabulary)
	return nil
}

func (v *vocabularyRepository) UpsertOne(ctx context.Context, id string, vocabulary *vocabulary_domain.Vocabulary) (vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	filterReference := bson.M{"_id": vocabulary.UnitID}
	count, err := collectionUnit.CountDocuments(ctx, filterReference)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}

	if count == 0 {
		return vocabulary_domain.Response{}, errors.New("the lesson ID do not exist")
	}

	doc, err := internal.ToDoc(vocabulary)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collectionVocabulary.FindOneAndUpdate(ctx, query, update, opts)

	var updatePost vocabulary_domain.Response
	if err := res.Decode(&updatePost); err != nil {
		return vocabulary_domain.Response{}, err
	}

	return updatePost, nil
}

func (v *vocabularyRepository) DeleteOne(ctx context.Context, vocabularyID string) error {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionMean := v.database.Collection(v.collectionMean)
	collectionMark := v.database.Collection(v.collectionMark)

	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}

	filterChild := bson.M{
		"vocabulary_id": objID,
	}
	countChildMean, err := collectionMean.CountDocuments(ctx, filterChild)
	if err != nil {
		return err
	}
	if countChildMean > 0 {
		return errors.New(`lesson cannot remove`)
	}

	countChildMark, err := collectionMark.CountDocuments(ctx, filterChild)
	if err != nil {
		return err
	}
	if countChildMark > 0 {
		return errors.New(`lesson cannot remove`)
	}

	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`the lesson is removed`)
	}

	_, err = collectionVocabulary.DeleteOne(ctx, filter)
	return err
}
