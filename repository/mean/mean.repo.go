package mean

import (
	mean_domain "clean-architecture/domain/mean"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type meanRepository struct {
	database             mongo.Database
	collectionMean       string
	collectionVocabulary string
}

func NewMeanRepository(db mongo.Database, collectionMean string, collectionVocabulary string) mean_domain.IMeanRepository {
	return &meanRepository{
		database:             db,
		collectionMean:       collectionMean,
		collectionVocabulary: collectionVocabulary,
	}
}

func (m *meanRepository) FetchMany(ctx context.Context) ([]mean_domain.Response, error) {
	collectionMean := m.database.Collection(m.collectionMean)
	collectionVocabulary := m.database.Collection(m.collectionVocabulary)

	cursor, err := collectionMean.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer func(cursor mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var means []mean_domain.Response
	for cursor.Next(ctx) {
		var mean mean_domain.Mean
		var mean2 mean_domain.Response
		if err := cursor.Decode(&mean); err != nil {
			return nil, err
		}

		var vocabulary vocabulary_domain.Vocabulary
		err := collectionVocabulary.FindOne(ctx, bson.M{"_id": mean.VocabularyID}).Decode(&vocabulary)
		if err != nil {
			return nil, err
		}

		// Gắn tên của course vào lesson
		mean2.VocabularyID = vocabulary.Id

		// Thêm lesson vào slice lessons
		means = append(means, mean2)
	}
	err = cursor.All(ctx, &means)
	if means == nil {
		return []mean_domain.Response{}, err
	}
	return means, err
}

func (m *meanRepository) CreateOne(ctx context.Context, mean *mean_domain.Mean, fieldOfIT string) error {
	collectionMean := m.database.Collection(m.collectionMean)
	collectionVocabulary := m.database.Collection(m.collectionVocabulary)

	filterMean := bson.M{"name": mean.Description, "vocabulary_id": mean.VocabularyID}
	filterVocabulary := bson.M{"_id": mean.VocabularyID}
	filterCondition := bson.M{"field_of_it": fieldOfIT}

	// check exists with CountDocuments
	countVocabulary, err := collectionVocabulary.CountDocuments(ctx, filterVocabulary)
	if err != nil {
		return err
	}

	countMean, err := collectionMean.CountDocuments(ctx, filterMean)
	if err != nil {
		return err
	}

	countCond, err := collectionVocabulary.CountDocuments(ctx, filterCondition)
	if countMean > 0 {
		return errors.New("the unit name in lesson did exist")
	}
	if countVocabulary == 0 {
		return errors.New("the lesson ID do not exist")
	}
	if countCond > 0 && countMean > 0 {
		return errors.New("the means is unique in one lesson")
	}
	_, err = collectionMean.InsertOne(ctx, mean)
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
		return errors.New(`the unit is removed`)
	}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}
