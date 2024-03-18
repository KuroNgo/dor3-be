package unit_repo

import (
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type unitRepository struct {
	database         mongo.Database
	collectionUnit   string
	collectionLesson string
}

func NewUnitRepository(db mongo.Database, collectionUnit string, collectionLesson string) unit_domain.IUnitRepository {
	return &unitRepository{
		database:         db,
		collectionUnit:   collectionUnit,
		collectionLesson: collectionLesson,
	}
}

func (u *unitRepository) FetchMany(ctx context.Context) ([]unit_domain.Response, error) {
	collectionUnit := u.database.Collection(u.collectionUnit)
	collectionLesson := u.database.Collection(u.collectionLesson)

	cursor, err := collectionUnit.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var units []unit_domain.Response
	for cursor.Next(ctx) {
		var unit unit_domain.Unit
		var unit2 unit_domain.Response
		if err := cursor.Decode(&unit); err != nil {
			return nil, err
		}

		var lesson lesson_domain.Lesson
		err := collectionLesson.FindOne(ctx, bson.M{"_id": unit.LessonID}).Decode(&lesson)
		if err != nil {
			return nil, err
		}

		// Gắn tên của course vào lesson
		unit2.LessonID = lesson.ID

		// Thêm lesson vào slice lessons
		units = append(units, unit2)
	}
	err = cursor.All(ctx, &units)
	if units == nil {
		return []unit_domain.Response{}, err
	}
	return units, err
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
	collection := u.database.Collection(u.collectionLesson)
	objID, err := primitive.ObjectIDFromHex(unitID)
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