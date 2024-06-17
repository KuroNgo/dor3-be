package unit_repo

import (
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	"clean-architecture/internal"
	"clean-architecture/internal/cache"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"
)

type unitRepository struct {
	database              *mongo.Database
	collectionUnit        string
	collectionUnitProcess string
	collectionLesson      string
	collectionVocabulary  string
	collectionExam        string
	collectionExercise    string
	collectionQuiz        string
}

// NewUnitRepository hàm khởi tạo (constructor) để khởi tạo instance của struct
func NewUnitRepository(db *mongo.Database, collectionUnit string, collectionUnitProcess string, collectionLesson string, collectionVocabulary string, collectionExam string, collectionExercise string, collectionQuiz string) unit_domain.IUnitRepository {
	return &unitRepository{
		database:              db,
		collectionUnit:        collectionUnit,
		collectionUnitProcess: collectionUnitProcess,
		collectionLesson:      collectionLesson,
		collectionVocabulary:  collectionVocabulary,
		collectionExam:        collectionExam,
		collectionExercise:    collectionExercise,
		collectionQuiz:        collectionQuiz,
	}
}

var (
	unitsCache            = cache.NewTTL[string, []unit_domain.UnitResponse]()
	unitCache             = cache.NewTTL[string, unit_domain.UnitResponse]()
	unitPrimOIDCache      = cache.NewTTL[string, primitive.ObjectID]()
	unitsUserProcessCache = cache.NewTTL[string, []unit_domain.UnitProcessResponse]()
	unitUserProcessCache  = cache.NewTTL[string, unit_domain.UnitProcessResponse]()
	detailUnitCache       = cache.NewTTL[string, unit_domain.DetailResponse]()

	wg           sync.WaitGroup
	mu           sync.Mutex
	rx           sync.RWMutex
	isProcessing bool
)

// FetchManyInUser fetches a paginated list of unit processes for a given user.
// It attempts to retrieve data from cache first; if not found, it fetches from the database.
// The function returns a list of unit processes, detail response with pagination info, and an error if any.
func (u *unitRepository) FetchManyInUser(ctx context.Context, user primitive.ObjectID, page string) ([]unit_domain.UnitProcessResponse, unit_domain.DetailResponse, error) {
	// Create a buffered error channel to handle errors from goroutines
	errCh := make(chan error, 1)
	// Create channels to retrieve unit processes and detail responses
	unitsUserProcessCh := make(chan []unit_domain.UnitProcessResponse)
	detailCh := make(chan unit_domain.DetailResponse)

	// Increment wait group counter by 2 for the two goroutines
	wg.Add(2)

	// Goroutine to fetch unit processes from cache
	go func() {
		defer wg.Done() // Decrement the wait group counter when the goroutine completes
		data, found := unitsUserProcessCache.Get(user.Hex() + page)
		if found {
			unitsUserProcessCh <- data // Send the data to the channel if found
		}
	}()

	// Goroutine to fetch detail from cache
	go func() {
		defer wg.Done() // Decrement the wait group counter when the goroutine completes
		data, found := detailUnitCache.Get(user.Hex() + "detail")
		if found {
			detailCh <- data // Send the data to the channel if found
		}
	}()

	// Goroutine to close the channels after wait group completes
	go func() {
		defer close(detailCh)
		defer close(unitsUserProcessCh)
		wg.Wait() // Wait for both cache fetch goroutines to complete
	}()

	// Read data from the channels
	unitsUserProcessData := <-unitsUserProcessCh
	detailData := <-detailCh
	if !internal.IsZeroValue(unitsUserProcessData) && !internal.IsZeroValue(detailData) {
		return unitsUserProcessData, detailData, nil // Return data if both are non-zero
	}

	// MongoDB collections
	collectionUnit := u.database.Collection(u.collectionUnit)
	collectionUnitProcess := u.database.Collection(u.collectionUnitProcess)

	// Filter for user's unit processes
	filterUnitProcessByUser := bson.M{"user_id": user}

	// Pagination setup
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	// Count total lessons
	countUnit, err := collectionUnit.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	// Calculate total pages
	totalPages := (countUnit + int64(perPage) - 1) / int64(perPage)

	// Count user's unit processes
	count, err := collectionUnitProcess.CountDocuments(ctx, filterUnitProcessByUser)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	var unitsProcess []unit_domain.UnitProcessResponse
	if count < countUnit {
		cursorUnit, err := collectionUnit.Find(ctx, bson.M{})
		if err != nil {
			return nil, unit_domain.DetailResponse{}, err
		}
		defer func(cursorUnit *mongo.Cursor, ctx context.Context) {
			err := cursorUnit.Close(ctx)
			if err != nil {
				errCh <- err
				return
			}
		}(cursorUnit, ctx)

		for cursorUnit.Next(ctx) {
			var unit unit_domain.Unit
			if err = cursorUnit.Decode(&unit); err != nil {
				return nil, unit_domain.DetailResponse{}, err
			}
			wg.Add(1)
			go func(unit unit_domain.Unit) {
				defer wg.Done()
				unitProcess := unit_domain.UnitProcess{
					UnitID:     unit.ID,
					LessonID:   unit.LessonID,
					UserID:     user,
					IsComplete: 0,
				}

				filter := bson.M{"unit_id": unit.ID}
				countUnitChild, err := collectionUnit.CountDocuments(ctx, filter)
				if err != nil {
					errCh <- err
					return
				}

				if countUnitChild == 0 {
					_, err := collectionUnitProcess.InsertOne(ctx, &unitProcess)
					if err != nil {
						errCh <- err
						return
					}
				}
			}(unit)
		}
		wg.Wait()

		cursor, err := collectionUnitProcess.Find(ctx, filterUnitProcessByUser, findOptions)
		if err != nil {
			return nil, unit_domain.DetailResponse{}, err
		}
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			err := cursor.Close(ctx)
			if err != nil {
				errCh <- err
				return
			}
		}(cursor, ctx)
	}

	cursor, err := collectionUnitProcess.Find(ctx, filterUnitProcessByUser, findOptions)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var unitProcess unit_domain.UnitProcess
		if err := cursor.Decode(&unitProcess); err != nil {
			return nil, unit_domain.DetailResponse{}, err
		}

		wg.Add(1)
		go func(unitProcess unit_domain.UnitProcess) {
			defer wg.Done()
			filter := bson.M{"_id": unitProcess.UnitID}
			var unit unit_domain.Unit
			err = collectionUnit.FindOne(ctx, filter).Decode(&unit)

			var unitProcessRes = unit_domain.UnitProcessResponse{
				Unit:               unit,
				UserID:             user,
				IsComplete:         unitProcess.IsComplete,
				ExamIsComplete:     unitProcess.ExamIsComplete,
				ExerciseIsComplete: unitProcess.ExerciseIsComplete,
				QuizIsComplete:     unitProcess.QuizIsComplete,
				TotalScore:         unitProcess.TotalScore,
			}

			mu.Lock()
			unitsProcess = append(unitsProcess, unitProcessRes)
			mu.Unlock()
		}(unitProcess)
	}
	wg.Wait()

	if err := cursor.Err(); err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}
	sort.Sort(unit_domain.UnitProcessResponseList(unitsProcess))

	detail := unit_domain.DetailResponse{
		Page:        totalPages,
		CurrentPage: pageNumber,
	}

	select {
	case err = <-errCh:
		return nil, unit_domain.DetailResponse{}, err
	default:
		return unitsProcess, detail, nil
	}
}

func (u *unitRepository) FetchOneByIDInUser(ctx context.Context, user primitive.ObjectID, id string) (unit_domain.UnitProcessResponse, error) {
	unitCh := make(chan unit_domain.UnitProcessResponse, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := unitUserProcessCache.Get(user.Hex() + id)
		if found {
			unitCh <- data
			return
		}
	}()

	go func() {
		defer close(unitCh)
		wg.Wait()
	}()

	unitData := <-unitCh
	if !internal.IsZeroValue(unitData) {
		return unitData, nil
	}

	collectionUnit := u.database.Collection(u.collectionUnit)
	collectionUnitProcess := u.database.Collection(u.collectionUnitProcess)

	idUnit, _ := primitive.ObjectIDFromHex(id)
	filterUnitProcessByUser := bson.M{"unit_id": idUnit, "user_id": user}
	var unitProcess unit_domain.UnitProcess
	err := collectionUnitProcess.FindOne(ctx, filterUnitProcessByUser).Decode(&unitProcess)
	if err != nil {
		return unit_domain.UnitProcessResponse{}, err
	}

	filter := bson.M{"_id": idUnit}
	var unit unit_domain.Unit
	err = collectionUnit.FindOne(ctx, filter).Decode(&unit)
	if err != nil {
		return unit_domain.UnitProcessResponse{}, err
	}

	unitProcessRes := unit_domain.UnitProcessResponse{
		Unit:               unit,
		UserID:             user,
		IsComplete:         unitProcess.IsComplete,
		ExamIsComplete:     unitProcess.ExamIsComplete,
		ExerciseIsComplete: unitProcess.ExerciseIsComplete,
		QuizIsComplete:     unitProcess.QuizIsComplete,
		TotalScore:         unitProcess.TotalScore,
	}

	unitUserProcessCache.Set(user.Hex()+id, unitProcessRes, 5*time.Minute)
	return unitProcessRes, nil

}

func (u *unitRepository) FetchManyNotPaginationInUser(ctx context.Context, user primitive.ObjectID) ([]unit_domain.UnitProcessResponse, error) {
	// Create a buffered error channel to handle errors from goroutines
	errCh := make(chan error, 1)

	// Create channels to retrieve unit processes and detail responses
	unitsUserProcessCh := make(chan []unit_domain.UnitProcessResponse, 1)

	// Increment wait group counter by 2 for the two goroutines
	wg.Add(2)

	// Goroutine to fetch unit processes from cache
	go func() {
		defer wg.Done() // Decrement the wait group counter when the goroutine completes
		data, found := unitsUserProcessCache.Get(user.Hex())
		if found {
			unitsUserProcessCh <- data // Send the data to the channel if found
		}
	}()

	// Goroutine to close the channels after wait group completes
	go func() {
		defer close(unitsUserProcessCh)
		wg.Wait() // Wait for both cache fetch goroutines to complete
	}()

	// Read data from the channels
	unitsUserProcessData := <-unitsUserProcessCh
	if !internal.IsZeroValue(unitsUserProcessData) {
		return unitsUserProcessData, nil // Return data if both are non-zero
	}

	// MongoDB collections
	collectionUnit := u.database.Collection(u.collectionUnit)
	collectionUnitProcess := u.database.Collection(u.collectionUnitProcess)

	// Filter for user's unit processes
	filterUnitProcessByUser := bson.M{"user_id": user}

	// Count total lessons
	countUnit, err := collectionUnit.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	// Count user's unit processes
	count, err := collectionUnitProcess.CountDocuments(ctx, filterUnitProcessByUser)
	if err != nil {
		return nil, err
	}

	var unitsProcess []unit_domain.UnitProcessResponse
	if count < countUnit {
		cursorUnit, err := collectionUnit.Find(ctx, bson.M{})
		if err != nil {
			return nil, err
		}
		defer func(cursorUnit *mongo.Cursor, ctx context.Context) {
			err := cursorUnit.Close(ctx)
			if err != nil {
				errCh <- err
				return
			}
		}(cursorUnit, ctx)

		for cursorUnit.Next(ctx) {
			var unit unit_domain.Unit
			if err = cursorUnit.Decode(&unit); err != nil {
				return nil, err
			}

			wg.Add(1)
			go func(unit unit_domain.Unit) {
				defer wg.Done()
				unitProcess := unit_domain.UnitProcess{
					UnitID:     unit.ID,
					LessonID:   unit.LessonID,
					UserID:     user,
					IsComplete: 0,
				}

				filter := bson.M{"unit_id": unit.ID}
				countUnitChild, err := collectionUnit.CountDocuments(ctx, filter)
				if err != nil {
					errCh <- err
					return
				}

				if countUnitChild == 0 {
					_, err = collectionUnitProcess.InsertOne(ctx, &unitProcess)
					if err != nil {
						errCh <- err
						return
					}
				}
			}(unit)
		}
		wg.Wait()

		cursor, err := collectionUnitProcess.Find(ctx, filterUnitProcessByUser)
		if err != nil {
			return nil, err
		}
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			err := cursor.Close(ctx)
			if err != nil {
				errCh <- err
				return
			}
		}(cursor, ctx)
	}

	cursor, err := collectionUnitProcess.Find(ctx, filterUnitProcessByUser)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var unitProcess unit_domain.UnitProcess
		if err := cursor.Decode(&unitProcess); err != nil {
			return nil, err
		}

		wg.Add(1)
		go func(unitProcess unit_domain.UnitProcess) {
			defer wg.Done()
			filter := bson.M{"_id": unitProcess.UnitID}
			var unit unit_domain.Unit
			err = collectionUnit.FindOne(ctx, filter).Decode(&unit)

			var unitProcessRes = unit_domain.UnitProcessResponse{
				Unit:               unit,
				UserID:             user,
				IsComplete:         unitProcess.IsComplete,
				ExamIsComplete:     unitProcess.ExamIsComplete,
				ExerciseIsComplete: unitProcess.ExerciseIsComplete,
				QuizIsComplete:     unitProcess.QuizIsComplete,
				TotalScore:         unitProcess.TotalScore,
			}

			mu.Lock()
			unitsProcess = append(unitsProcess, unitProcessRes)
			mu.Unlock()
		}(unitProcess)
	}
	wg.Wait()

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	sort.Sort(unit_domain.UnitProcessResponseList(unitsProcess))

	select {
	case err = <-errCh:
		return nil, err
	default:
		return unitsProcess, nil
	}
}

func (u *unitRepository) FetchByIdLessonInUser(ctx context.Context, user primitive.ObjectID, idLesson string, page string) ([]unit_domain.UnitProcessResponse, unit_domain.DetailResponse, error) {
	errCh := make(chan error)

	unitsUserProcessCh := make(chan []unit_domain.UnitProcessResponse)
	detailCh := make(chan unit_domain.DetailResponse)

	wg.Add(2)
	go func() {
		defer wg.Done()
		data, found := unitsUserProcessCache.Get(user.Hex() + idLesson + page)
		if found {
			unitsUserProcessCh <- data
		}
	}()

	go func() {
		defer wg.Done()
		data, found := detailUnitCache.Get(user.Hex() + idLesson + "detail")
		if found {
			detailCh <- data
		}
	}()

	go func() {
		defer close(detailCh)
		defer close(unitsUserProcessCh)
		wg.Wait()
	}()

	unitsUserProcessData := <-unitsUserProcessCh
	detailData := <-detailCh
	if !internal.IsZeroValue(unitsUserProcessData) && !internal.IsZeroValue(detailData) {
		return unitsUserProcessData, detailData, nil
	}

	collectionUnit := u.database.Collection(u.collectionUnit)
	collectionUnitProcess := u.database.Collection(u.collectionUnitProcess)

	// Thực hiện phân trang
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip)).SetSort(bson.M{"level": 1})

	lessonID, _ := primitive.ObjectIDFromHex(idLesson)

	// Đếm số lượng khóa học trong collection 'lesson'
	filterUnit := bson.M{"lesson_id": lessonID}
	countUnit, err := collectionUnit.CountDocuments(ctx, filterUnit)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	// Đếm số lượng CourseProcess của người dùng
	filterLessonProcessByUser := bson.M{"lesson_id": lessonID, "user_id": user}
	count, err := collectionUnitProcess.CountDocuments(ctx, filterLessonProcessByUser)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}

	// Tính toán tổng số trang dựa trên số lượng khóa học và số khóa học mỗi trang
	totalPages := (count + int64(perPage) - 1) / int64(perPage)

	var unitsProcess []unit_domain.UnitProcessResponse
	// Nếu không có LessonProcess cho người dùng, khởi tạo chúng
	if count != countUnit {
		cursorUnit, err := collectionUnit.Find(ctx, bson.D{})
		if err != nil {
			return nil, unit_domain.DetailResponse{}, err
		}
		defer func(cursorUnit *mongo.Cursor, ctx context.Context) {
			err = cursorUnit.Close(ctx)
			if err != nil {
				return
			}
		}(cursorUnit, ctx)

		for cursorUnit.Next(ctx) {
			var unit unit_domain.Unit
			if err = cursorUnit.Decode(&unit); err != nil {
				return nil, unit_domain.DetailResponse{}, err
			}

			wg.Add(1)
			go func(unit unit_domain.Unit) {
				defer wg.Done()
				unitProcess := unit_domain.UnitProcess{
					LessonID:   unit.LessonID,
					UnitID:     unit.ID,
					UserID:     user,
					IsComplete: 0,
				}

				// Thực hiện tìm kiếm theo name để kiểm tra có dữ liệu trùng không
				filter := bson.M{"unit_id": unit.ID}
				countUnitChild, err := collectionUnit.CountDocuments(ctx, filter)
				if err != nil {
					errCh <- err
					return
				}

				if countUnitChild == 0 {
					_, err = collectionUnitProcess.InsertOne(ctx, &unitProcess)
					if err != nil {
						log.Println("Error inserting course process:", err)
						errCh <- err
						return
					}
				}
			}(unit)
		}
		wg.Wait()

		// Tìm các LessonProcess của người dùng với phân trang
		cursor, err := collectionUnitProcess.Find(ctx, filterLessonProcessByUser, findOptions)
		if err != nil {
			return nil, unit_domain.DetailResponse{}, err
		}
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			err := cursor.Close(ctx)
			if err != nil {
				errCh <- err
				return
			}
		}(cursor, ctx)
	}

	// Tìm các LessonProcess của người dùng với phân trang
	cursor, err := collectionUnitProcess.Find(ctx, filterLessonProcessByUser, findOptions)
	if err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	// Đọc dữ liệu từ cursor và thêm vào slice LessonProcess
	for cursor.Next(ctx) {
		var unitProcess unit_domain.UnitProcess
		if err := cursor.Decode(&unitProcess); err != nil {
			return nil, unit_domain.DetailResponse{}, err
		}

		wg.Add(1)
		go func(unitProcess unit_domain.UnitProcess) {
			defer wg.Done()
			filter := bson.M{"_id": unitProcess.UnitID}
			var unit unit_domain.Unit
			err = collectionUnit.FindOne(ctx, filter).Decode(&unit)

			var unitProcessRes = unit_domain.UnitProcessResponse{
				Unit:               unit,
				UserID:             unitProcess.UserID,
				IsComplete:         unitProcess.IsComplete,
				ExamIsComplete:     unitProcess.ExamIsComplete,
				ExerciseIsComplete: unitProcess.ExerciseIsComplete,
				QuizIsComplete:     unitProcess.QuizIsComplete,
				TotalScore:         unitProcess.TotalScore,
			}

			mu.Lock()
			unitsProcess = append(unitsProcess, unitProcessRes)
			mu.Unlock()
		}(unitProcess)
	}
	wg.Wait()

	if err := cursor.Err(); err != nil {
		return nil, unit_domain.DetailResponse{}, err
	}
	sort.Sort(unit_domain.UnitProcessResponseList(unitsProcess))

	// Lấy thống kê cho detail response
	detail := unit_domain.DetailResponse{
		Page:        totalPages,
		CurrentPage: pageNumber,
	}

	unitsUserProcessCache.Set(user.Hex()+idLesson+page, unitsProcess, 5*time.Minute)
	detailUnitCache.Set(user.Hex()+idLesson+"detail", detail, 5*time.Minute)

	select {
	case err = <-errCh:
		return nil, unit_domain.DetailResponse{}, err
	default:
		return unitsProcess, detail, nil
	}
}

func (u *unitRepository) UpdateCompleteInUser(ctx context.Context, user primitive.ObjectID) (*mongo.UpdateResult, error) {
	//TODO implement me
	panic("implement me")
}

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
	// Channel to save units
	unitsCh := make(chan []unit_domain.UnitResponse, 1)
	// Channel to save detail
	detailCh := make(chan unit_domain.DetailResponse, 1)

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
		data, found := detailUnitCache.Get("detail")
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
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip)).SetSort(bson.M{"level": 1})

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
	for cursor.Next(ctx) {
		var unit unit_domain.UnitResponse
		if err = cursor.Decode(&unit); err != nil {
			return nil, unit_domain.DetailResponse{}, err
		}

		wg.Add(1)
		go func(unit unit_domain.UnitResponse) {
			defer wg.Done()
			// Fetch related data for each unit
			countVocabulary, err := u.countVocabularyByUnitID(ctx, unit.ID)
			if err != nil {
				errCh <- err
				return
			}

			// Set fetched data to unit
			unit.CountVocabulary = countVocabulary
			mu.Lock()
			units = append(units, unit)
			mu.Unlock()
		}(unit)
	}
	wg.Wait()

	// Prepare detail response
	detail := unit_domain.DetailResponse{
		Page:        totalPages,
		CurrentPage: pageNumber,
	}

	// Cache the results
	unitsCache.Set(page, units, 5*time.Minute)
	detailUnitCache.Set("detail", detail, 5*time.Minute)

	// Return results or error
	select {
	case err = <-errCh:
		return nil, unit_domain.DetailResponse{}, err
	default:
		return units, detail, nil
	}
}

func (u *unitRepository) FetchManyNotPaginationInAdmin(ctx context.Context) ([]unit_domain.UnitResponse, error) {
	errCh := make(chan error, 1)
	// Channel to save units
	unitsCh := make(chan []unit_domain.UnitResponse, 1)
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

	for cursor.Next(ctx) {
		var unit unit_domain.UnitResponse
		if err = cursor.Decode(&unit); err != nil {
			errCh <- err
			return nil, err
		}

		wg.Add(1)
		go func(unit unit_domain.UnitResponse) {
			defer wg.Done()
			countVocabulary, err := u.countVocabularyByUnitID(ctx, unit.ID)
			if err != nil {
				errCh <- err
				return
			}

			unit.CountVocabulary = countVocabulary
			units = append(units, unit)
		}(unit)
	}
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
	errCh := make(chan error, 1)
	unitsCh := make(chan []unit_domain.UnitResponse)
	detailCh := make(chan unit_domain.DetailResponse)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := unitsCache.Get(idLesson + page)
		if found {
			unitsCh <- data
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := detailUnitCache.Get(idLesson + "detail")
		if found {
			detailCh <- data
		}
	}()

	go func() {
		defer close(unitsCh)
		defer close(detailCh)
		wg.Wait()
	}()

	unitsData := <-unitsCh
	detailData := <-detailCh
	if !internal.IsZeroValue(unitsData) && !internal.IsZeroValue(detailData) {
		return unitsData, detailData, nil
	}

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
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var units []unit_domain.UnitResponse
	for cursor.Next(ctx) {
		var unit unit_domain.UnitResponse
		if err = cursor.Decode(&unit); err != nil {
			return nil, unit_domain.DetailResponse{}, err
		}

		wg.Add(1)
		go func(unit unit_domain.UnitResponse) {
			defer wg.Done()
			countVocabulary, err := u.countVocabularyByUnitID(ctx, unit.ID)
			if err != nil {
				errCh <- err
				return
			}

			// Gắn LessonID vào đơn vị
			unit.LessonID = idLesson2
			unit.CountVocabulary = countVocabulary

			units = append(units, unit)
		}(unit)
	}
	wg.Wait()

	response := unit_domain.DetailResponse{
		Page:        totalPages,
		CurrentPage: pageNumber,
	}

	unitsCache.Set(idLesson+page, units, 5*time.Minute)
	detailUnitCache.Set(idLesson+"detail", response, 5*time.Minute)

	select {
	case err = <-errCh:
		return nil, unit_domain.DetailResponse{}, err
	default:
		return units, response, nil
	}
}

func (u *unitRepository) FetchOneByIDInAdmin(ctx context.Context, id string) (unit_domain.UnitResponse, error) {
	unitCh := make(chan unit_domain.UnitResponse)
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := unitCache.Get(id)
		if found {
			unitCh <- data
		}
	}()

	go func() {
		defer close(unitCh)
		wg.Wait()
	}()

	unitData := <-unitCh
	if !internal.IsZeroValue(unitData) {
		return unitData, nil
	}

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

	unitCache.Set(id, unit, 5*time.Minute)
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

	// Clear data value in cache memory
	wg.Add(2)
	go func() {
		defer wg.Done()
		unitsCache.Clear()
	}()

	// clear data value in cache memory due to increase num
	go func() {
		defer wg.Done()
		detailUnitCache.Clear()
	}()

	wg.Wait()
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
		if err != nil {
			return err
		}

		// Clear data value in cache memory
		wg.Add(2)
		go func() {
			defer wg.Done()
			unitsCache.Clear()
		}()

		// clear data value in cache memory due to increase num
		go func() {
			defer wg.Done()
			detailUnitCache.Clear()
		}()

		wg.Wait()
		return nil
	}

	return errors.New("the unit cannot be created because the vocabulary in the latest unit is not complete")
}

func (u *unitRepository) UpdateOneInAdmin(ctx context.Context, unit *unit_domain.Unit) (*mongo.UpdateResult, error) {
	// Khóa lock giúp bảo vệ course
	mu.Lock()
	defer mu.Unlock()

	if isProcessing {
		return nil, errors.New("another goroutine is already processing")
	}

	isProcessing = true
	defer func() {
		isProcessing = false
	}()

	collection := u.database.Collection(u.collectionUnit)

	filter := bson.D{{Key: "_id", Value: unit.ID}}
	update := bson.M{"$set": unit}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Clear data value in cache memory
	wg.Add(3)
	go func() {
		defer wg.Done()
		unitsCache.Clear()
	}()

	go func() {
		defer wg.Done()
		unitCache.Remove(unit.ID.Hex())
	}()

	// clear data value in cache memory due to increase num
	go func() {
		defer wg.Done()
		detailUnitCache.Clear()
	}()

	wg.Wait()
	return data, nil
}

func (u *unitRepository) FindUnitIDByUnitLevelInAdmin(ctx context.Context, unitLevel int, fieldOfIT string) (primitive.ObjectID, error) {
	unitPrimOIDCh := make(chan primitive.ObjectID)
	go func() {
		data, found := unitPrimOIDCache.Get(strconv.Itoa(unitLevel) + fieldOfIT)
		if found {
			unitPrimOIDCh <- data
		}
	}()

	go func() {
		defer close(unitPrimOIDCh)
		wg.Wait()
	}()

	unitData := <-unitPrimOIDCh
	if !internal.IsZeroValue(unitData) {
		return unitData, nil
	}
	collectionUnit := u.database.Collection(u.collectionUnit)
	collectionLesson := u.database.Collection(u.collectionLesson)

	// Tìm lesson
	var lessons []lesson_domain.Lesson
	cursor, err := collectionLesson.Find(ctx, bson.D{})
	for cursor.Next(ctx) {
		var lesson lesson_domain.Lesson
		if err := cursor.Decode(&lesson); err != nil {
			return primitive.NilObjectID, err
		}

		lessons = append(lessons, lesson)
	}

	var unitMain unit_domain.Unit
	for _, data := range lessons {
		if fieldOfIT == data.Name {
			var lesson lesson_domain.Lesson

			filterLesson := bson.M{"name": fieldOfIT}
			err = collectionLesson.FindOne(ctx, filterLesson).Decode(&lesson)
			if err != nil {
				return primitive.NilObjectID, err
			}

			var unit unit_domain.Unit
			filterUnit := bson.M{"lesson_id": lesson.ID, "level": unitLevel}
			err = collectionUnit.FindOne(ctx, filterUnit).Decode(&unit)
			if err != nil {
				return primitive.NilObjectID, err
			}

			unitMain = unit
			break
		}
	}

	unitPrimOIDCache.Set(strconv.Itoa(unitLevel)+fieldOfIT, unitMain.ID, 5*time.Minute)
	return unitMain.ID, nil
}

func (u *unitRepository) DeleteOneInAdmin(ctx context.Context, unitID string) error {
	// Khóa lock giúp bảo vệ course
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
	if err != nil {
		return err
	}
	// Clear data value in cache memory
	wg.Add(3)
	go func() {
		defer wg.Done()
		unitsCache.Clear()
	}()

	go func() {
		defer wg.Done()
		unitCache.Remove(unitID)
	}()

	// clear data value in cache memory due to increase num
	go func() {
		defer wg.Done()
		detailUnitCache.Clear()
	}()
	wg.Wait()

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
