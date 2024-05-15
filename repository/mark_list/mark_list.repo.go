package mark_list_repository

import (
	mark_list_domain "clean-architecture/domain/mark_list"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type markListRepository struct {
	database                 *mongo.Database
	collectionMarkList       string
	collectionMarkVocabulary string
}

func NewListRepository(db *mongo.Database, collectionMarkList string, collectionMarkVocabulary string) mark_list_domain.IMarkListRepository {
	return &markListRepository{
		database:                 db,
		collectionMarkList:       collectionMarkList,
		collectionMarkVocabulary: collectionMarkVocabulary,
	}
}

func (m *markListRepository) FetchManyByUserID(ctx context.Context, userID string) (mark_list_domain.Response, error) {
	collectionMarkList := m.database.Collection(m.collectionMarkList)
	//collectionMarkVocabulary := m.database.Collection(m.collectionMarkVocabulary)

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return mark_list_domain.Response{}, err
	}

	filter := bson.M{"user_id": idUser}
	cursor, err := collectionMarkList.Find(ctx, filter)
	if err != nil {
		return mark_list_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	//count, err := collectionMarkList.CountDocuments(ctx, bson.M{})
	//if err != nil {
	//	return mark_list_domain.Response{}, err
	//}
	//countVocabulary, err := collectionMarkVocabulary.CountDocuments(ctx, bson.M{})
	//if err != nil {
	//	return mark_list_domain.Response{}, err
	//}

	var markLists []mark_list_domain.MarkList

	for cursor.Next(ctx) {
		var markList mark_list_domain.MarkList
		if err = cursor.Decode(&markList); err != nil {
			return mark_list_domain.Response{}, err
		}

		// Gắn CourseID vào bài học
		markList.UserID = idUser

		markLists = append(markLists, markList)
	}

	response := mark_list_domain.Response{
		MarkList: markLists,
	}

	return response, nil
}

func (m *markListRepository) FetchById(ctx context.Context, id string) (mark_list_domain.MarkList, error) {
	collectionMarkList := m.database.Collection(m.collectionMarkList)

	idMarkList, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mark_list_domain.MarkList{}, err
	}

	filter := bson.M{"_id": idMarkList}
	var markList mark_list_domain.MarkList
	err = collectionMarkList.FindOne(ctx, filter).Decode(&markList)
	if err != nil {
		return mark_list_domain.MarkList{}, err
	}

	return markList, err
}

func (m *markListRepository) FetchMany(ctx context.Context) (mark_list_domain.Response, error) {
	collectionMarkList := m.database.Collection(m.collectionMarkList)

	cursor, err := collectionMarkList.Find(ctx, bson.D{})
	if err != nil {
		return mark_list_domain.Response{}, err
	}

	var markLists []mark_list_domain.MarkList
	for cursor.Next(ctx) {
		var markList mark_list_domain.MarkList
		if err = cursor.Decode(&markList); err != nil {
			return mark_list_domain.Response{}, err
		}

		// Thêm lesson vào slice lessons
		markLists = append(markLists, markList)
	}

	statisticsCh := make(chan mark_list_domain.Statistics)
	go func() {
		statistics, _ := m.Statistics(ctx)
		statisticsCh <- statistics
	}()
	statistics := <-statisticsCh

	markListRes := mark_list_domain.Response{
		MarkList:   markLists,
		Statistics: statistics,
	}

	return markListRes, err
}

func (m *markListRepository) UpdateOne(ctx context.Context, markList *mark_list_domain.MarkList) (*mongo.UpdateResult, error) {
	collectionMarkList := m.database.Collection(m.collectionMarkList)

	filter := bson.M{"_id": markList.ID, "user_id": markList.UserID}
	update := bson.M{
		"$set": bson.M{
			"name_list":   markList.NameList,
			"description": markList.Description,
		},
	}

	data, err := collectionMarkList.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (m *markListRepository) CreateOne(ctx context.Context, markList *mark_list_domain.MarkList) error {
	collectionMarkList := m.database.Collection(m.collectionMarkList)

	filter := bson.M{"name_list": markList.NameList}
	// check exists with CountDocuments
	count, err := collectionMarkList.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the mark list name did exist")
	}

	_, err = collectionMarkList.InsertOne(ctx, markList)
	return err
}

func (m *markListRepository) UpsertOne(ctx context.Context, id string, markList *mark_list_domain.MarkList) (mark_list_domain.Response, error) {
	collectionMarkList := m.database.Collection(m.collectionMarkList)

	doc, err := internal.ToDoc(markList)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mark_list_domain.Response{}, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collectionMarkList.FindOneAndUpdate(ctx, query, update, opts)

	var updatedPost mark_list_domain.Response

	if err := res.Decode(&updatedPost); err != nil {
		return mark_list_domain.Response{}, errors.New("no post with that Id exists")
	}

	return updatedPost, nil
}

func (m *markListRepository) DeleteOne(ctx context.Context, markListID string) error {
	collectionMarkList := m.database.Collection(m.collectionMarkList)
	collectionMarkVocabulary := m.database.Collection(m.collectionMarkList)

	// Convert courseID string to ObjectID
	objID, err := primitive.ObjectIDFromHex(markListID)
	if err != nil {
		return err
	}

	filterChildren := bson.M{
		"mark_list_id": markListID,
	}

	//Check if any lesson is associated with the course
	countFK, err := collectionMarkVocabulary.CountDocuments(ctx, filterChildren)
	if err != nil {
		return err
	}
	if countFK > 0 {
		return errors.New("the mark list cannot be deleted because it is associated with mark vocabulary")
	}

	// Delete the mark list
	filter := bson.M{"_id": objID}
	result, err := collectionMarkList.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result == nil {
		return errors.New("the mark list was not found or already deleted")
	}

	return nil
}

// countLessonsByCourseID counts the number of lessons associated with a course.
func (m *markListRepository) countMarkVocabularyByMarkListID(ctx context.Context, courseID string) (int64, error) {
	collectionMarkList := m.database.Collection(m.collectionMarkList)

	filter := bson.M{"course_id": courseID}
	count, err := collectionMarkList.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *markListRepository) Statistics(ctx context.Context) (mark_list_domain.Statistics, error) {
	collectionMarkList := m.database.Collection(m.collectionMarkList)
	collectionMarkVocabulary := m.database.Collection(m.collectionMarkVocabulary)

	countMarkList, err := collectionMarkList.CountDocuments(ctx, bson.D{})
	if err != nil {
		return mark_list_domain.Statistics{}, err
	}

	countMarkVocabulary, err := collectionMarkVocabulary.CountDocuments(ctx, bson.D{})
	if err != nil {
		return mark_list_domain.Statistics{}, err
	}

	statistics := mark_list_domain.Statistics{
		Total:           countMarkList,
		CountVocabulary: countMarkVocabulary,
	}
	return statistics, nil
}
