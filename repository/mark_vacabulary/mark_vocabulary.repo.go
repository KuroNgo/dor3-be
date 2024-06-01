package mark_vacabulary_repository

import (
	mark_list_domain "clean-architecture/domain/mark_list"
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type markVocabularyRepository struct {
	database                 *mongo.Database
	collectionMarkList       string
	collectionVocabulary     string
	collectionMarkVocabulary string
}

func NewMarkVocabularyRepository(db *mongo.Database, collectionMarkList string, collectionVocabulary string, collectionMarkVocabulary string) mark_vocabulary_domain.IMarkToFavouriteRepository {
	return &markVocabularyRepository{
		database:                 db,
		collectionMarkList:       collectionMarkList,
		collectionVocabulary:     collectionVocabulary,
		collectionMarkVocabulary: collectionMarkVocabulary,
	}
}

var (
	wg sync.WaitGroup
)

func (m *markVocabularyRepository) FetchManyByMarkListID(ctx context.Context, markListId string) ([]mark_vocabulary_domain.MarkToFavourite, error) {
	collectionMarkVocabulary := m.database.Collection(m.collectionMarkList)

	// Tạo các bộ lọc
	filterMarkList := bson.D{{Key: "mark_list_id", Value: markListId}}
	_, err := collectionMarkVocabulary.Find(ctx, filterMarkList)
	if err != nil {
		return nil, err
	}

	cursor, err := collectionMarkVocabulary.Find(ctx, filterMarkList)
	if err != nil {
		return nil, err
	}

	var markVocabs []mark_vocabulary_domain.MarkToFavourite
	for cursor.Next(ctx) {
		var markVocab mark_vocabulary_domain.MarkToFavourite
		if err = cursor.Decode(&markVocab); err != nil {
			return nil, err
		}

		markVocabs = append(markVocabs, markVocab)
	}

	return markVocabs, nil
}

func (m *markVocabularyRepository) FetchManyByMarkListIDAndUserId(ctx context.Context, markListID string, userID string) (mark_vocabulary_domain.Response, error) {
	collectionMarkVocabulary := m.database.Collection(m.collectionMarkVocabulary)
	collectionVocabulary := m.database.Collection(m.collectionVocabulary)
	collectionMarkList := m.database.Collection(m.collectionMarkList)

	// Chuyển đổi markListID và userID sang ObjectID
	markListObjID, err := primitive.ObjectIDFromHex(markListID)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}

	// Tạo các bộ lọc
	filterMarkList := bson.D{{Key: "mark_list_id", Value: markListObjID}}

	// Tìm các mark vocabulary của mark list
	cursor, err := collectionMarkVocabulary.Find(ctx, filterMarkList)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var markVocabularies []mark_vocabulary_domain.MarkToFavouriteResponse
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var markVocabulary mark_vocabulary_domain.MarkToFavourite
			if err = cursor.Decode(&markVocabulary); err != nil {
				return
			}

			var markVocabularyRes mark_vocabulary_domain.MarkToFavouriteResponse
			if err = cursor.Decode(&markVocabulary); err != nil {
				return
			}

			// Gắn MarkListID vào mark vocabulary
			var vocabulary vocabulary_domain.Vocabulary
			filterVocabulary := bson.M{"_id": markVocabulary.VocabularyID}
			_ = collectionVocabulary.FindOne(ctx, filterVocabulary).Decode(&vocabulary)

			var markList mark_list_domain.MarkList
			filterMarkList := bson.D{{Key: "_id", Value: markListObjID}}
			_ = collectionMarkList.FindOne(ctx, filterMarkList).Decode(&markList)

			markVocabularyRes.ID = markVocabulary.ID
			markVocabularyRes.UserId = markVocabulary.UserId
			markVocabularyRes.Vocabulary = vocabulary
			markVocabularyRes.MarkList = markList

			markVocabularies = append(markVocabularies, markVocabularyRes)
		}
	}()

	wg.Wait()

	// Tạo và trả về response
	response := mark_vocabulary_domain.Response{
		Total:                   len(markVocabularies),
		MarkToFavouriteResponse: markVocabularies,
	}
	return response, nil
}

func (m *markVocabularyRepository) FetchManyByMarkList(ctx context.Context, markListId string) (mark_vocabulary_domain.Response, error) {
	collectionMarkList := m.database.Collection(m.collectionMarkList)
	collectionMarkVocabulary := m.database.Collection(m.collectionMarkVocabulary)
	collectionVocabulary := m.database.Collection(m.collectionVocabulary)

	idMarkList, err := primitive.ObjectIDFromHex(markListId)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}

	filterParent := bson.M{"_id": idMarkList}
	countParent, err := collectionMarkList.CountDocuments(ctx, filterParent)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}
	if countParent == 0 {
		return mark_vocabulary_domain.Response{}, errors.New("the mark_list_id not found")
	}

	filter := bson.M{"mark_list_id": idMarkList}
	cursor, err := collectionMarkVocabulary.Find(ctx, filter)
	if err != nil {
		return mark_vocabulary_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var markVocabularies []mark_vocabulary_domain.MarkToFavouriteResponse

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var markVocabulary mark_vocabulary_domain.MarkToFavourite
			if err = cursor.Decode(&markVocabulary); err != nil {
				return
			}

			// Gắn MarkListID vào mark vocabulary
			var vocabulary vocabulary_domain.Vocabulary
			filterVocabulary := bson.M{"_id": markVocabulary.VocabularyID}
			err = collectionVocabulary.FindOne(ctx, filterVocabulary).Decode(&vocabulary)
			if err != nil {
				return
			}

			var markList mark_list_domain.MarkList
			err = collectionMarkList.FindOne(ctx, filter).Decode(&markList)
			if err != nil {
				return
			}

			var markVocabularyRes mark_vocabulary_domain.MarkToFavouriteResponse
			if err = cursor.Decode(&markVocabulary); err != nil {
				return
			}
			markVocabularyRes.Vocabulary = vocabulary
			markVocabularyRes.MarkList = markList

			markVocabularies = append(markVocabularies, markVocabularyRes)
		}
	}()

	wg.Wait()

	response := mark_vocabulary_domain.Response{
		Total:                   len(markVocabularies),
		MarkToFavouriteResponse: markVocabularies,
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
	// Lấy các collection cụ thể từ cơ sở dữ liệu MongoDB
	collectionMarkVocabulary := m.database.Collection(m.collectionMarkVocabulary)
	collectionMarkList := m.database.Collection(m.collectionMarkList)
	collectionVocabulary := m.database.Collection(m.collectionVocabulary)

	// Kiểm tra xem Mark List có tồn tại không bằng cách đếm các tài liệu khớp với MarkListID
	filterMarkListReference := bson.M{"_id": markVocabulary.MarkListID}
	countMarkListParent, err := collectionMarkList.CountDocuments(ctx, filterMarkListReference)
	if err != nil {
		return err // Trả về lỗi nếu có lỗi xảy ra khi đếm tài liệu
	}
	if countMarkListParent == 0 {
		return errors.New("the mark list does not exist") // Trả về lỗi nếu Mark List không tồn tại
	}

	// Kiểm tra xem Vocabulary có tồn tại không bằng cách đếm các tài liệu khớp với VocabularyID
	filterMarkVocabularyReference := bson.M{"_id": markVocabulary.VocabularyID}
	countMarkVocabularyParent, err := collectionVocabulary.CountDocuments(ctx, filterMarkVocabularyReference)
	if err != nil {
		return err // Trả về lỗi nếu có lỗi xảy ra khi đếm tài liệu
	}
	if countMarkVocabularyParent == 0 {
		return errors.New("the vocabulary does not exist") // Trả về lỗi nếu Vocabulary không tồn tại
	}

	// Kiểm tra xem Mark Vocabulary đã tồn tại chưa bằng cách đếm các tài liệu khớp với MarkListID và VocabularyID
	filter := bson.M{"mark_list_id": markVocabulary.MarkListID, "vocabulary_id": markVocabulary.VocabularyID}
	count, err := collectionMarkVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err // Trả về lỗi nếu có lỗi xảy ra khi đếm tài liệu
	}
	if count > 0 {
		return errors.New("the mark vocabulary already exists") // Trả về lỗi nếu Mark Vocabulary đã tồn tại
	}

	// Thực hiện chèn dữ liệu mới vào collection MarkVocabulary
	_, err = collectionMarkVocabulary.InsertOne(ctx, markVocabulary)
	return err // Trả về kết quả của lệnh chèn hoặc lỗi nếu có lỗi xảy ra
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
