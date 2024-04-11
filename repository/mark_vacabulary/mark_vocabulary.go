package mark_vacabulary_repository

import (
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type markVocabularyRepository struct {
	database                 mongo.Database
	collectionMarkList       string
	collectionVocabulary     string
	collectionMarkVocabulary string
}

func NewMarkVocabularyRepository(db mongo.Database, collectionMarkList string, collectionVocabulary string, collectionMarkVocabulary string) mark_vocabulary_domain.IMarkToFavouriteRepository {
	return &markVocabularyRepository{
		database:                 db,
		collectionMarkList:       collectionMarkList,
		collectionVocabulary:     collectionVocabulary,
		collectionMarkVocabulary: collectionMarkVocabulary,
	}
}

func (m *markVocabularyRepository) FetchManyByMarkListIDAndUserId(ctx context.Context, markListID string, userID string) (mark_vocabulary_domain.Response, error) {
	collectionMarkVocabulary := m.database.Collection(m.collectionMarkVocabulary)

	// Chuyển đổi markListID và userID sang ObjectID
	markListObjID, err := primitive.ObjectIDFromHex(markListID)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}
	userIDObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}

	// Tạo các bộ lọc
	filterMarkList := bson.D{{Key: "mark_list_id", Value: markListObjID}}
	filterUser := bson.D{{Key: "user_id", Value: userIDObjID}}

	// Đếm số lượng mark vocabulary
	countMarkList, err := collectionMarkVocabulary.CountDocuments(ctx, filterMarkList)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}

	// Đếm số lượng mark vocabulary của người dùng
	countUser, err := collectionMarkVocabulary.CountDocuments(ctx, filterUser)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}

	// Tìm các mark vocabulary của mark list
	cursor, err := collectionMarkVocabulary.Find(ctx, filterMarkList)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var markVocabularies []mark_vocabulary_domain.MarkToFavourite
	for cursor.Next(ctx) {
		var markVocabulary mark_vocabulary_domain.MarkToFavourite
		if err := cursor.Decode(&markVocabulary); err != nil {
			return mark_vocabulary_domain.Response{}, err
		}
		// Gắn MarkListID vào mark vocabulary
		markVocabulary.MarkListID = markListObjID
		markVocabularies = append(markVocabularies, markVocabulary)
	}

	// Tạo và trả về response
	response := mark_vocabulary_domain.Response{
		MarkToFavourite: markVocabularies,
		CountMarkList:   countMarkList,
		CountUser:       countUser,
	}
	return response, nil
}

func (m *markVocabularyRepository) UpdateOne(ctx context.Context, markVocabularyID string, markVocabulary mark_vocabulary_domain.MarkToFavourite) error {
	collection := m.database.Collection(m.collectionMarkVocabulary)
	objID, err := primitive.ObjectIDFromHex(markVocabularyID)
	doc, err := internal.ToDoc(markVocabulary)

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": doc}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (m *markVocabularyRepository) CreateOne(ctx context.Context, markVocabulary *mark_vocabulary_domain.MarkToFavourite) error {
	collectionMarkVocabulary := m.database.Collection(m.collectionMarkVocabulary)
	collectionMarkList := m.database.Collection(m.collectionMarkList)
	collectionVocabulary := m.database.Collection(m.collectionVocabulary)

	filter := bson.M{"mark_list_id": markVocabulary.MarkListID, "mark_vocabulary_id": markVocabulary.MarkListID}

	// check exists with CountDocuments
	count, err := collectionMarkVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the vocabulary in mark list did exist")
	}

	filterMarkListReference := bson.M{"_id": markVocabulary.MarkListID}
	countMarkListParent, err := collectionMarkList.CountDocuments(ctx, filterMarkListReference)
	if err != nil {
		return err
	}
	if countMarkListParent == 0 {
		return errors.New("the mark list do not exists")
	}

	filterMarkVocabularyReference := bson.M{"_id": markVocabulary.VocabularyID}
	countMarkVocabularyParent, err := collectionVocabulary.CountDocuments(ctx, filterMarkVocabularyReference)
	if err != nil {
		return err
	}
	if countMarkVocabularyParent == 0 {
		return errors.New("the vocabulary do not exists")
	}

	_, err = collectionMarkVocabulary.InsertOne(ctx, markVocabulary)
	return err
}

func (m *markVocabularyRepository) DeleteOne(ctx context.Context, markVocabularyID string) error {
	collectionMarkVocabulary := m.database.Collection(m.collectionMarkVocabulary)
	objID, err := primitive.ObjectIDFromHex(markVocabularyID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}

	count, err := collectionMarkVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`the mark vocabulary is removed or have not exist`)
	}

	_, err = collectionMarkVocabulary.DeleteOne(ctx, filter)
	return err
}
