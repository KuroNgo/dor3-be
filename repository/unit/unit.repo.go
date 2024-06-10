package unit_repo

import (
	unit_domain "clean-architecture/domain/unit"
	"clean-architecture/internal"
	"clean-architecture/internal/cache"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"sync"
	"time"
)

type unitRepository struct {
	database             *mongo.Database
	collectionUnit       string
	collectionLesson     string
	collectionVocabulary string
	collectionExam       string
	collectionExercise   string
	collectionQuiz       string
}

// NewUnitRepository hàm khởi tạo (constructor) để khởi tạo instance của struct
func NewUnitRepository(db *mongo.Database, collectionUnit string, collectionLesson string, collectionVocabulary string, collectionExam string, collectionExercise string, collectionQuiz string) unit_domain.IUnitRepository {
	return &unitRepository{
		database:             db,
		collectionUnit:       collectionUnit,
		collectionLesson:     collectionLesson,
		collectionVocabulary: collectionVocabulary,
		collectionExam:       collectionExam,
		collectionExercise:   collectionExercise,
		collectionQuiz:       collectionQuiz,
	}
}

var (
	unitsCache  = cache.NewTTL[string, []unit_domain.UnitResponse]()
	unitCache   = cache.NewTTL[string, unit_domain.UnitResponse]()
	detailCache = cache.NewTTL[string, unit_domain.DetailResponse]()

	wg           sync.WaitGroup
	mu           sync.Mutex
	isProcessing bool
)

// FetchManyInAdmin retrieves multiple unit responses and a detail response for a given page number.
// It first attempts to retrieve data from the cache. If the data is not found in the cache,
// it queries the database to fetch the data. The results are then cached for future use.
//
// Parameters:
// - ctx: context.Context: The context to control cancellations and timeouts.
// - page: string: The page number as a string, used for pagination.
//
// Returns:
// - []unit_domain.UnitResponse: A slice of UnitResponse containing the fetched unit data.
// - unit_domain.DetailResponse: A DetailResponse containing pagination details.
// - error: An error object if an error occurred during the process, or nil if successful.
func (u *unitRepository) FetchManyInAdmin(ctx context.Context, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	// Channel to log errors
	errCh := make(chan error, 1)
	defer close(errCh)
	// Channel to save units
	unitsCh := make(chan []unit_domain.UnitResponse)
	// Channel to save detail
	detailCh := make(chan unit_domain.DetailResponse)

	// Use WaitGroup to wait for all goroutines to complete
	wg.Add(2)

	// Goroutine to fetch units from cache
	go func() {
		defer wg.Done()
		data, found := unitsCache.Get(page)
		if found {
			unitsCh <- data
			return
		}
	}()

	// Goroutine to fetch detail from cache
	go func() {
		defer wg.Done()
		data, found := detailCache.Get("detail")
		if found {
			detailCh <- data
			return
		}
	}()

	// Goroutine to close channels once WaitGroup is done
	go func() {
		defer close(unitsCh)
		defer close(detailCh)
		wg.Wait()
	}()

	// Read from channels
	unitsData := <-unitsCh
	detailData := <-detailCh
	if !internal.IsZeroValue(unitsData) && !internal.IsZeroValue(detailData) {
		return unitsData, detailData, nil
	}

	// MongoDB collections
	collectionUnit := u.database.Collection(u.collectionUnit)

	// Pagination
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 5
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	// Count total units
	count, err := collectionUnit.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	// Calculate total pages
	totalPages := (count + int64(perPage) - 1) / int64(perPage)

	// Find units with pagination
	cursor, err := collectionUnit.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var units []unit_domain.UnitResponse

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var unit unit_domain.UnitResponse
			if err = cursor.Decode(&unit); err != nil {
				errCh <- err
				return
			}

			// Fetch related data for each unit
			countVocabulary, err := u.countVocabularyByUnitID(ctx, unit.ID)
			if err != nil {
				errCh <- err
				return
			}

			// Set fetched data to unit
			unit.CountVocabulary = countVocabulary
			units = append(units, unit)
		}
	}()
	wg.Wait()

	// Prepare detail response
	detail := unit_domain.DetailResponse{
		Page:        totalPages,
		CurrentPage: pageNumber,
	}

	// Cache the results
	unitsCache.Set(page, units, 5*time.Minute)
	detailCache.Set("detail", detail, 5*time.Minute)

	// Return results or error
	select {
	case err = <-errCh:
		return nil, unit_domain.DetailResponse{}, err
	default:
		return units, detail, nil
	}
}

func (u *unitRepository) FetchManyNotPaginationInAdmin(ctx context.Context) ([]unit_domain.UnitResponse, error) {
	errCh := make(chan error)
	defer close(errCh)
	// Channel to save units
	unitsCh := make(chan []unit_domain.UnitResponse)
	wg.Add(1)
	// Goroutine to fetch units from cache
	go func() {
		defer wg.Done()
		data, found := unitsCache.Get("units")
		if found {
			unitsCh <- data
			return
		}
	}()

	// Goroutine to close channels once WaitGroup is done
	go func() {
		defer close(unitsCh)
		wg.Wait()
	}()

	// Read from channels
	unitsData := <-unitsCh
	if !internal.IsZeroValue(unitsData) {
		return unitsData, nil
	}

	collectionUnit := u.database.Collection(u.collectionUnit)

	cursor, err := collectionUnit.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var units []unit_domain.UnitResponse
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var unit unit_domain.UnitResponse
			if err = cursor.Decode(&unit); err != nil {
				errCh <- err
				return
			}

			countVocabulary, err := u.countVocabularyByUnitID(ctx, unit.ID)
			if err != nil {
				errCh <- err
				return
			}

			unit.CountVocabulary = countVocabulary
			units = append(units, unit)
		}

	}()
	wg.Wait()

	// Cache the results
	unitsCache.Set("units", units, 5*time.Minute)

	// Return results or error
	select {
	case err = <-errCh:
		return nil, err
	default:
		return units, nil
	}
}

func (u *unitRepository) FetchByIdLessonInAdmin(ctx context.Context, idLesson string, page string) ([]unit_domain.UnitResponse, unit_domain.DetailResponse, error) {
	collectionUnit := u.database.Collection(u.collectionUnit)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 5
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip)).SetSort(bson.D{{"level", 1}})

	count, err := collectionUnit.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	// Calculate total pages
	totalPages := (count + int64(perPage) - 1) / int64(perPage)

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

		countVocabulary, err := u.countVocabularyByUnitID(ctx, unit.ID)
		if err != nil {
			return nil, unit_domain.DetailResponse{}, err
		}

		// Gắn LessonID vào đơn vị
		unit.LessonID = idLesson2
		unit.CountVocabulary = countVocabulary

		units = append(units, unit)
	}

	response := unit_domain.DetailResponse{
		Page:        totalPages,
		CurrentPage: pageNumber,
	}
	return units, response, nil
}

func (u *unitRepository) FetchOneByIDInAdmin(ctx context.Context, id string) (unit_domain.UnitResponse, error) {
	collectionUnit := u.database.Collection(u.collectionUnit)

	idUnit, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return unit_domain.UnitResponse{}, err
	}

	filter := bson.M{"_id": idUnit}
	var unit unit_domain.UnitResponse
	err = collectionUnit.FindOne(ctx, filter).Decode(&unit)
	if err != nil {
		return unit_domain.UnitResponse{}, err
	}

	countVocabulary, err := u.countVocabularyByUnitID(ctx, unit.ID)
	if err != nil {
		return unit_domain.UnitResponse{}, err
	}

	unit.CountVocabulary = countVocabulary
	return unit, nil
}

func (u *unitRepository) CreateOneByNameLessonInAdmin(ctx context.Context, unit *unit_domain.Unit) error {
	mu.Lock()
	defer mu.Unlock()

	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	isProcessing = true
	defer func() {
		isProcessing = false
	}()

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

func (u *unitRepository) CreateOneInAdmin(ctx context.Context, unit *unit_domain.Unit) error {
	mu.Lock()
	defer mu.Unlock()

	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	isProcessing = true
	defer func() {
		isProcessing = false
	}()

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
		return nil
	}

	return errors.New("the unit cannot be created because the vocabulary in the latest unit is not complete")
}

func (u *unitRepository) UpdateOneInAdmin(ctx context.Context, unit *unit_domain.Unit) (*mongo.UpdateResult, error) {
	collection := u.database.Collection(u.collectionUnit)

	filter := bson.D{{Key: "_id", Value: unit.ID}}
	update := bson.M{"$set": unit}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *unitRepository) UpdateCompleteInUser(ctx context.Context) (*mongo.UpdateResult, error) {
	//TODO implement me
	panic("implement me")
}

func (u *unitRepository) DeleteOneInAdmin(ctx context.Context, unitID string) error {
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

// countVocabularyByUnitID counts the number of lessons associated with a course.
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
