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
	database           *mongo.Database
	collectionMarkList string
}

func NewListRepository(db *mongo.Database, collectionMarkList string) mark_list_domain.IMarkListRepository {
	return &markListRepository{
		database:           db,
		collectionMarkList: collectionMarkList,
	}
}

func (m *markListRepository) FetchManyByUserID(ctx context.Context, userID string) (mark_list_domain.Response, error) {
	collectionMarkList := m.database.Collection(m.collectionMarkList)

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return mark_list_domain.Response{}, err
	}

	filter := bson.M{"user_id": idUser}
	cursor, err := collectionMarkList.Find(ctx, filter)
	if err != nil {
		return mark_list_domain.Response{}, err
	}
	defer cursor.Close(ctx)

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
	err = cursor.All(ctx, &markLists)
	courseRes := mark_list_domain.Response{
		MarkList: markLists,
	}

	return courseRes, err
}

func (m *markListRepository) UpdateOne(ctx context.Context, markListID string, markList mark_list_domain.MarkList) error {
	collectionMarkList := m.database.Collection(m.collectionMarkList)

	doc, err := internal.ToDoc(markList)
	objID, err := primitive.ObjectIDFromHex(markListID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collectionMarkList.UpdateOne(ctx, filter, update)

	return err
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
