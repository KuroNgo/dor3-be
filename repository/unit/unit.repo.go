package unit_repo

import (
	unit_domain "clean-architecture/domain/unit"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type unitRepository struct {
	database             mongo.Database
	collectionUnit       string
	collectionLesson     string
	collectionVocabulary string
}

func NewUnitRepository(db mongo.Database, collectionUnit string, collectionLesson string, collectionVocabulary string) unit_domain.IUnitRepository {
	return &unitRepository{
		database:             db,
		collectionUnit:       collectionUnit,
		collectionLesson:     collectionLesson,
		collectionVocabulary: collectionVocabulary,
	}
}
func (u *unitRepository) CreateOneByNameLesson(ctx context.Context, unit *unit_domain.Unit) error {
	collectionUnit := u.database.Collection(u.collectionUnit)
	collectionLesson := u.database.Collection(u.collectionLesson)

	filter := bson.M{"name": unit.Name, "lesson_id": unit.LessonID}

	filterParent := bson.M{"_id": unit.LessonID}
	countParent, err := collectionLesson.CountDocuments(ctx, filterParent)
	if err != nil {
		return err
	}
	if countParent == 0 {
		return errors.New("parent lesson not found")
	}

	count, err := collectionUnit.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the unit name already exists in the lesson")
	}

	_, err = collectionUnit.InsertOne(ctx, unit)
	if err != nil {
		return err
	}
	return nil
}

func (u *unitRepository) FindLessonIDByLessonName(ctx context.Context, lessonName string) (primitive.ObjectID, error) {
	collectionLesson := u.database.Collection(u.collectionLesson)

	filter := bson.M{"name": lessonName}
	var data struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	err := collectionLesson.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return data.Id, nil
}

func (u *unitRepository) FetchByIdLesson(ctx context.Context, idLesson string) (unit_domain.Response, error) {
	collectionUnit := u.database.Collection(u.collectionUnit)

	idLesson2, err := primitive.ObjectIDFromHex(idLesson)
	if err != nil {
		return unit_domain.Response{}, err
	}

	filter := bson.M{"lesson_id": idLesson2}

	cursor, err := collectionUnit.Find(ctx, filter)
	if err != nil {
		return unit_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var units []unit_domain.Unit

	for cursor.Next(ctx) {
		var unit unit_domain.Unit
		if err := cursor.Decode(&unit); err != nil {
			return unit_domain.Response{}, err
		}

		// Gắn LessonID vào đơn vị
		unit.LessonID = idLesson2

		units = append(units, unit)
	}

	response := unit_domain.Response{
		Unit: units,
	}
	return response, nil
}

func (u *unitRepository) UpdateComplete(ctx context.Context, updateData unit_domain.Update) error {
	collection := u.database.Collection(u.collectionUnit)
	objID, err := primitive.ObjectIDFromHex(updateData.UnitID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "is_complete", Value: updateData.IsComplete},
		{Key: "who_updates", Value: updateData.WhoUpdate},
	}}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	isLessonComplete, err := u.CheckLessonComplete(ctx, updateData.LessonID)
	if err != nil {
		return err
	}

	lessonCollection := u.database.Collection(u.collectionLesson)
	lessonObjID, err := primitive.ObjectIDFromHex(updateData.LessonID)
	if err != nil {
		return err
	}

	lessonUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "is_complete", Value: isLessonComplete}}}}
	lessonFilter := bson.D{{Key: "_id", Value: lessonObjID}}
	_, err = lessonCollection.UpdateOne(ctx, lessonFilter, lessonUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (u *unitRepository) CheckLessonComplete(ctx context.Context, lessonID string) (bool, error) {
	collection := u.database.Collection(u.collectionUnit)

	lessonID2, _ := primitive.ObjectIDFromHex(lessonID)

	cursor, err := collection.Find(ctx, bson.D{{Key: "lesson_id", Value: lessonID2}})
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return false, nil
	}

	for cursor.Next(ctx) {
		var unit unit_domain.Unit
		if err := cursor.Decode(&unit); err != nil {
			return false, err
		}

		if unit.IsComplete != 1 {
			return false, nil
		}
	}

	return true, nil
}

func (u *unitRepository) FetchMany(ctx context.Context, page string) (unit_domain.Response, error) {
	collectionUnit := u.database.Collection(u.collectionUnit)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return unit_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	cursor, err := collectionUnit.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return unit_domain.Response{}, err
	}
	defer func(cursor mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var units []unit_domain.Unit
	for cursor.Next(ctx) {
		var unit unit_domain.Unit
		if err := cursor.Decode(&unit); err != nil {
			return unit_domain.Response{}, err
		}

		// Thêm lesson vào slice lessons
		units = append(units, unit)
	}

	unitRes := unit_domain.Response{
		Unit: units,
	}

	return unitRes, err
}

func (u *unitRepository) CreateOne(ctx context.Context, unit *unit_domain.Unit) error {
	collectionUnit := u.database.Collection(u.collectionUnit)
	collectionLesson := u.database.Collection(u.collectionLesson)

	filterUnit := bson.M{"name": unit.Name, "lesson_id": unit.LessonID}
	filterLess := bson.M{"_id": unit.LessonID}

	// check exists with CountDocuments
	countLess, err := collectionLesson.CountDocuments(ctx, filterLess)
	if err != nil {
		return err
	}

	countUnit, err := collectionUnit.CountDocuments(ctx, filterUnit)
	if err != nil {
		return err
	}

	if countUnit > 0 {
		return errors.New("the unit name in lesson did exist")
	}
	if countLess == 0 {
		return errors.New("the lesson ID do not exist")
	}

	_, err = collectionUnit.InsertOne(ctx, unit)
	return nil
}

func (u *unitRepository) UpdateOne(ctx context.Context, unitID string, unit unit_domain.Unit) error {
	collection := u.database.Collection(u.collectionUnit)
	doc, err := internal.ToDoc(unit)
	objID, err := primitive.ObjectIDFromHex(unitID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (u *unitRepository) UpsertOne(ctx context.Context, id string, unit *unit_domain.Unit) (unit_domain.Response, error) {
	collectionUnit := u.database.Collection(u.collectionUnit)
	collectionLesson := u.database.Collection(u.collectionLesson)

	filterReference := bson.M{"_id": unit.LessonID}
	count, err := collectionLesson.CountDocuments(ctx, filterReference)
	if err != nil {
		return unit_domain.Response{}, err
	}

	if count == 0 {
		return unit_domain.Response{}, errors.New("the course ID do not exist")
	}

	doc, err := internal.ToDoc(unit)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return unit_domain.Response{}, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collectionUnit.FindOneAndUpdate(ctx, query, update, opts)

	var updatePost unit_domain.Response
	if err := res.Decode(&updatePost); err != nil {
		return unit_domain.Response{}, err
	}

	return updatePost, nil
}

func (u *unitRepository) DeleteOne(ctx context.Context, unitID string) error {
	collectionUnit := u.database.Collection(u.collectionUnit)
	collectionVocabulary := u.database.Collection(u.collectionVocabulary)
	objID, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}

	filterChild := bson.M{
		"unit_id": objID,
	}
	countChild, err := collectionVocabulary.CountDocuments(ctx, filterChild)
	if err != nil {
		return err
	}
	if countChild == 0 {
		return errors.New(`the unit can not remove`)
	}

	count, err := collectionUnit.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`the unit is removed`)
	}

	_, err = collectionUnit.DeleteOne(ctx, filter)
	return err
}
