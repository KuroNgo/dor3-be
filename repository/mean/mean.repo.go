package mean_repository

import (
	mean_domain "clean-architecture/domain/mean"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type meanRepository struct {
	database             mongo.Database
	collectionMean       string
	collectionVocabulary string
}

func (m *meanRepository) FindVocabularyIDByWord(ctx context.Context, unitName string) (primitive.ObjectID, error) {
	collectionVocabulary := m.database.Collection(m.collectionVocabulary)

	filter := bson.M{"name": unitName}
	var data struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	err := collectionVocabulary.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return data.Id, nil
}

func (m *meanRepository) CreateOneByWord(ctx context.Context, mean *mean_domain.Mean) error {
	collectionMean := m.database.Collection(m.collectionVocabulary)
	collectionVocabulary := m.database.Collection(m.collectionVocabulary)

	filter := bson.M{"description": mean.Description, "vocabulary_id": mean.VocabularyID}

	filterParent := bson.M{"_id": mean.VocabularyID}
	countParent, err := collectionVocabulary.CountDocuments(ctx, filterParent)
	if err != nil {
		return err
	}
	if countParent == 0 {
		return errors.New("parent unit not found")
	}

	count, err := collectionMean.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the mean already exists in the lesson")
	}

	_, err = collectionMean.InsertOne(ctx, mean)
	if err != nil {
		return err
	}
	return nil
}

func NewMeanRepository(db mongo.Database, collectionMean string, collectionVocabulary string) mean_domain.IMeanRepository {
	return &meanRepository{
		database:             db,
		collectionMean:       collectionMean,
		collectionVocabulary: collectionVocabulary,
	}
}

func (m *meanRepository) FetchMany(ctx context.Context) (mean_domain.Response, error) {
	collectionMean := m.database.Collection(m.collectionMean)
	collectionVocabulary := m.database.Collection(m.collectionVocabulary)

	cursor, err := collectionMean.Find(ctx, bson.D{})
	if err != nil {
		return mean_domain.Response{}, err
	}

	var means []mean_domain.Mean
	for cursor.Next(ctx) {
		var mean mean_domain.Mean
		if err := cursor.Decode(&mean); err != nil {
			return mean_domain.Response{}, err
		}

		var vocabulary vocabulary_domain.Vocabulary
		err = collectionVocabulary.FindOne(ctx, bson.M{"_id": mean.VocabularyID}).Decode(&vocabulary)
		if err != nil {
			return mean_domain.Response{}, err
		}

		// Gắn tên của course vào lesson
		mean.VocabularyID = vocabulary.Id

		// Thêm lesson vào slice lessons
		means = append(means, mean)
	}
	err = cursor.All(ctx, &means)
	meanRes := mean_domain.Response{
		Mean: means,
	}
	return meanRes, err
}

func (m *meanRepository) CreateOne(ctx context.Context, mean *mean_domain.Mean, fieldOfIT string) error {
	collectionMean := m.database.Collection(m.collectionMean)
	collectionVocabulary := m.database.Collection(m.collectionVocabulary)

	filterVocabulary := bson.M{"_id": mean.VocabularyID}
	countVocabulary, err := collectionVocabulary.CountDocuments(ctx, filterVocabulary)
	if err != nil {
		return err
	}
	if countVocabulary == 0 {
		return errors.New("the vocabulary does not exist")
	}

	filterMean := bson.M{"name": mean.Description, "vocabulary_id": mean.VocabularyID}
	countMean, err := collectionMean.CountDocuments(ctx, filterMean)
	if err != nil {
		return err
	}
	if countMean > 0 {
		return errors.New("the mean already exists for the vocabulary")
	}

	filterCondition := bson.M{"_id": mean.VocabularyID, "field_of_it": fieldOfIT}
	countCond, err := collectionVocabulary.CountDocuments(ctx, filterCondition)
	if err != nil {
		return err
	}
	if countCond == 0 {
		return errors.New("the vocabulary does not belong to the specified field of IT")
	}

	// Insert the mean into the database
	_, err = collectionMean.InsertOne(ctx, mean)
	if err != nil {
		return err
	}
	return nil
}

func (m *meanRepository) UpdateOne(ctx context.Context, meanID string, mean mean_domain.Mean) error {
	collection := m.database.Collection(m.collectionMean)
	doc, err := internal.ToDoc(mean)
	objID, err := primitive.ObjectIDFromHex(meanID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (m *meanRepository) UpsertOne(ctx context.Context, id string, mean *mean_domain.Mean) (mean_domain.Mean, error) {
	collectionMean := m.database.Collection(m.collectionMean)
	collectionVocabulary := m.database.Collection(m.collectionVocabulary)

	filterReference := bson.M{"_id": mean.VocabularyID}
	count, err := collectionVocabulary.CountDocuments(ctx, filterReference)
	if err != nil {
		return mean_domain.Mean{}, err
	}

	if count == 0 {
		return mean_domain.Mean{}, errors.New("the course ID do not exist")
	}

	doc, err := internal.ToDoc(mean)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mean_domain.Mean{}, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collectionMean.FindOneAndUpdate(ctx, query, update, opts)

	var updatePost mean_domain.Mean
	if err := res.Decode(&updatePost); err != nil {
		return mean_domain.Mean{}, err
	}

	return updatePost, nil
}

func (m *meanRepository) DeleteOne(ctx context.Context, meanID string) error {
	collection := m.database.Collection(m.collectionMean)
	objID, err := primitive.ObjectIDFromHex(meanID)
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
		return errors.New(`mean is removed or exist'`)
	}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}
