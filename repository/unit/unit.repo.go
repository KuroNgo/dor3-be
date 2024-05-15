package unit_repo

import (
	unit_domain "clean-architecture/domain/unit"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type unitRepository struct {
	database             *mongo.Database
	collectionUnit       string
	collectionLesson     string
	collectionVocabulary string
	errCh                chan error
}

var (
	units []unit_domain.UnitResponse
	unit  unit_domain.UnitResponse
)

func NewUnitRepository(db *mongo.Database, collectionUnit string, collectionLesson string, collectionVocabulary string) unit_domain.IUnitRepository {
	return &unitRepository{
		database:             db,
		collectionUnit:       collectionUnit,
		collectionLesson:     collectionLesson,
		collectionVocabulary: collectionVocabulary,
		errCh:                make(chan error),
	}
}

func (u *unitRepository) FetchMany(ctx context.Context, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	collectionUnit := u.database.Collection(u.collectionUnit)

	// pagination
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 5
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	// count unit
	calCh := make(chan int64)
	count, err := collectionUnit.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	go func() {
		defer close(calCh)
		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calCh <- cal1
		}
	}()

	cursor, err := collectionUnit.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	internal.Wg.Add(1)
	go func() {
		defer internal.Wg.Done()
		for cursor.Next(ctx) {
			if err = cursor.Decode(&unit); err != nil {
				return
			}

			countVocabulary, err := u.countVocabularyByUnitID(ctx, unit.ID)
			if err != nil {
				return
			}

			unit.CountVocabulary = countVocabulary
			units = append(units, unit)
		}

	}()
	internal.Wg.Wait()

	cal := <-calCh
	detail := unit_domain.DetailResponse{
		CountUnit:   count,
		Page:        cal,
		CurrentPage: pageNumber,
	}

	return units, detail, nil
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

func (u *unitRepository) FetchByIdLesson(ctx context.Context, idLesson string, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	collectionUnit := u.database.Collection(u.collectionUnit)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 5
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip)).SetSort(bson.D{{"level", 1}})

	calCh := make(chan int64)
	countUnitCh := make(chan int64)

	go func() {
		defer close(calCh)
		defer close(countUnitCh)
		count, err := collectionUnit.CountDocuments(ctx, bson.D{})
		if err != nil {
			return
		}

		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calCh <- cal1
		}
	}()

	idLesson2, err := primitive.ObjectIDFromHex(idLesson)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	filter := bson.M{"lesson_id": idLesson2}
	cursor, err := collectionUnit.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var units []unit_domain.UnitResponse
	for cursor.Next(ctx) {
		var unit unit_domain.UnitResponse
		if err := cursor.Decode(&unit); err != nil {
			return nil, unit_domain.DetailResponse{}, err
		}

		// Gắn LessonID vào đơn vị
		unit.LessonID = idLesson2

		units = append(units, unit)
	}

	cal := <-calCh
	countUnit := <-countUnitCh

	response := unit_domain.DetailResponse{
		CountUnit:   countUnit,
		Page:        cal,
		CurrentPage: pageNumber,
	}
	return units, response, nil
}

func (u *unitRepository) UpdateComplete(ctx context.Context, updateData *unit_domain.Unit) error {
	internal.Wg.Add(2)
	go func() {
		defer internal.Wg.Done()
		collection := u.database.Collection(u.collectionUnit)

		filter := bson.D{{Key: "_id", Value: updateData.ID}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "is_complete", Value: updateData.IsComplete},
			{Key: "learner", Value: updateData.Learner},
		}}}
		_, err := collection.UpdateOne(ctx, filter, &update)
		if err != nil {
			u.errCh <- err
			return
		}
	}()

	go func() {
		defer internal.Wg.Done()
		isLessonComplete, err := u.CheckLessonComplete(ctx, updateData.LessonID)
		if err != nil {
			u.errCh <- err
			return
		}

		lessonCollection := u.database.Collection(u.collectionLesson)
		if err != nil {
			u.errCh <- err
			return
		}

		lessonUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "is_complete", Value: isLessonComplete}}}}
		lessonFilter := bson.D{{Key: "_id", Value: updateData.LessonID}}
		_, err = lessonCollection.UpdateOne(ctx, lessonFilter, &lessonUpdate)
		if err != nil {
			u.errCh <- err
			return
		}
	}()

	internal.Wg.Wait()
	close(u.errCh)

	select {
	case err := <-u.errCh:
		return err
	default:
		return nil
	}
}

func (u *unitRepository) CheckLessonComplete(ctx context.Context, lessonID primitive.ObjectID) (bool, error) {
	collection := u.database.Collection(u.collectionUnit)

	cursor, err := collection.Find(ctx, bson.D{{Key: "lesson_id", Value: lessonID}})
	if err != nil {
		return false, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

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

// CreateOne: hệ thống sẽ tự tạo unit nếu số lượng vocabulary là 5
func (u *unitRepository) CreateOne(ctx context.Context, unit *unit_domain.Unit) error {
	collectionUnit := u.database.Collection(u.collectionUnit)
	collectionLesson := u.database.Collection(u.collectionLesson)
	collectionVocabulary := u.database.Collection(u.collectionVocabulary)

	filterUnit := bson.M{"name": unit.Name, "lesson_id": unit.LessonID}
	filterLess := bson.M{"_id": unit.LessonID}

	// check exists with CountDocuments
	countLess, err := collectionLesson.CountDocuments(ctx, filterLess)
	if err != nil {
		return err
	}
	if countLess == 0 {
		return errors.New("the lesson ID do not exist")
	}

	// đếm số lượng document trong unit
	countUnit, err := collectionUnit.CountDocuments(ctx, filterUnit)
	if err != nil {
		return err
	}
	if countUnit > 0 {
		return errors.New("the unit name in lesson did exist")
	}

	// tạo unit dựa trên vocabulary
	data, err := u.getLastUnit(ctx)
	filterVocabulary := bson.M{"unit_id": data.ID}
	countVocabulary, err := collectionVocabulary.CountDocuments(ctx, filterVocabulary)
	if err != nil {
		return err
	}
	if countVocabulary == 0 || countVocabulary > 5 {
		_, err = collectionUnit.InsertOne(ctx, unit)
	}

	return nil
}

func (u *unitRepository) UpdateOne(ctx context.Context, unit *unit_domain.Unit) (*mongo.UpdateResult, error) {
	collection := u.database.Collection(u.collectionUnit)

	filter := bson.D{{Key: "_id", Value: unit.ID}}
	update := bson.M{"$set": unit}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
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

// countLessonsByCourseID counts the number of lessons associated with a course.
func (u *unitRepository) countVocabularyByUnitID(ctx context.Context, unitID primitive.ObjectID) (int32, error) {
	collectionVocabulary := u.database.Collection(u.collectionVocabulary)

	filter := bson.M{"unit_id": unitID}
	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}

// getLastUnit lấy unit cuối cùng từ collection
func (u *unitRepository) getLastUnit(ctx context.Context) (*unit_domain.Unit, error) {
	collectionUnit := u.database.Collection(u.collectionUnit)
	findOptions := options.FindOne().SetSort(bson.D{{"_id", -1}})

	var unit unit_domain.Unit
	err := collectionUnit.FindOne(ctx, bson.D{}, findOptions).Decode(&unit)
	if err != nil {
		return nil, err
	}

	return &unit, nil
}
