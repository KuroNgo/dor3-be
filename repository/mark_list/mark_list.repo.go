package mark_list_repository

import (
	mark_list_domain "clean-architecture/domain/mark_list"
	"clean-architecture/internal"
	"clean-architecture/internal/cache/memory"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
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

var (
	markListResCache          = memory.NewTTL[string, mark_list_domain.Response]()
	markListCache             = memory.NewTTL[string, mark_list_domain.MarkList]()
	statisticsIndividualCache = memory.NewTTL[string, mark_list_domain.Statistics]()
	statisticsAdminCache      = memory.NewTTL[string, mark_list_domain.Statistics]()

	wg           sync.WaitGroup
	mu           sync.Mutex
	isProcessing bool
)

const (
	cacheTTL = 5 * time.Minute
)

func (m *markListRepository) FetchManyByUser(ctx context.Context, user primitive.ObjectID) (mark_list_domain.Response, error) {
	errCh := make(chan error, 1)
	markListResCh := make(chan mark_list_domain.Response, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := markListResCache.Get(user.Hex())
		if found {
			markListResCh <- data
		}
	}()

	go func() {
		defer close(markListResCh)
		wg.Wait()
	}()

	markListData := <-markListResCh
	if !internal.IsZeroValue(markListData) {
		return markListData, nil
	}

	collectionMarkList := m.database.Collection(m.collectionMarkList)

	filter := bson.M{"user_id": user}
	cursor, err := collectionMarkList.Find(ctx, filter)
	if err != nil {
		return mark_list_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var markLists []mark_list_domain.MarkList
	for cursor.Next(ctx) {
		var markList mark_list_domain.MarkList
		if err = cursor.Decode(&markList); err != nil {
			return mark_list_domain.Response{}, err
		}

		wg.Add(1)
		go func(markList mark_list_domain.MarkList) {
			defer wg.Done()
			// Gắn CourseID vào bài học
			markList.UserID = user
			markLists = append(markLists, markList)
		}(markList)
	}
	wg.Wait()

	var statistics mark_list_domain.Statistics
	go func() {
		statistics, _ = m.Statistics(ctx)
	}()

	response := mark_list_domain.Response{
		Statistics: statistics,
		MarkList:   markLists,
	}

	markListResCache.Set(user.Hex(), response, cacheTTL)

	select {
	case err = <-errCh:
		return mark_list_domain.Response{}, err
	default:
		return response, nil
	}
}

func (m *markListRepository) FetchByIdByUser(ctx context.Context, user primitive.ObjectID, id string) (mark_list_domain.MarkList, error) {
	markListCh := make(chan mark_list_domain.MarkList, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := markListCache.Get(user.Hex())
		if found {
			markListCh <- data
		}
	}()

	go func() {
		defer close(markListCh)
		wg.Wait()
	}()

	markListData := <-markListCh
	if !internal.IsZeroValue(markListData) {
		return markListData, nil
	}

	collectionMarkList := m.database.Collection(m.collectionMarkList)

	idMarkList, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mark_list_domain.MarkList{}, err
	}

	filter := bson.M{"_id": idMarkList, "user_id": user}
	var markList mark_list_domain.MarkList
	err = collectionMarkList.FindOne(ctx, filter).Decode(&markList)
	if err != nil {
		return mark_list_domain.MarkList{}, err
	}

	markListCache.Set(user.Hex(), markList, cacheTTL)
	return markList, err
}

func (m *markListRepository) UpdateOneByUser(ctx context.Context, user primitive.ObjectID, markList *mark_list_domain.MarkList) (*mongo.UpdateResult, error) {
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

	wg.Add(3)
	go func() {
		defer wg.Done()
		markListCache.Clear()
	}()

	go func() {
		defer wg.Done()
		statisticsIndividualCache.Clear()
	}()

	go func() {
		defer wg.Done()
		statisticsAdminCache.Clear()
	}()
	wg.Wait()

	return data, err
}

func (m *markListRepository) CreateOneByUser(ctx context.Context, user primitive.ObjectID, markList *mark_list_domain.MarkList) error {
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

	collectionMarkList := m.database.Collection(m.collectionMarkList)

	filter := bson.M{"name_list": markList.NameList, "user_id": user}
	// check exists with CountDocuments
	count, err := collectionMarkList.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the mark list name did exist")
	}

	_, err = collectionMarkList.InsertOne(ctx, markList)
	if err != nil {
		return err
	}

	wg.Add(3)
	go func() {
		defer wg.Done()
		markListResCache.Clear()
	}()

	go func() {
		defer wg.Done()
		statisticsIndividualCache.Clear()
	}()

	go func() {
		defer wg.Done()
		statisticsAdminCache.Clear()
	}()
	wg.Wait()

	return err
}

func (m *markListRepository) UpsertOneByUser(ctx context.Context, user primitive.ObjectID, id string, markList *mark_list_domain.MarkList) (mark_list_domain.Response, error) {
	// Khóa lock giúp bảo vệ course
	mu.Lock()
	defer mu.Unlock()

	if isProcessing {
		return mark_list_domain.Response{}, errors.New("another goroutine is already processing")
	}

	isProcessing = true
	defer func() {
		isProcessing = false
	}()

	collectionMarkList := m.database.Collection(m.collectionMarkList)

	doc, err := internal.ToDoc(markList)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mark_list_domain.Response{}, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{"_id", idHex}, {"user_id", user}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collectionMarkList.FindOneAndUpdate(ctx, query, update, opts)

	var updatedPost mark_list_domain.Response
	if err = res.Decode(&updatedPost); err != nil {
		return mark_list_domain.Response{}, errors.New("no post with that Id exists")
	}

	wg.Add(4)
	go func() {
		defer wg.Done()
		markListResCache.Clear()
	}()

	go func() {
		defer wg.Done()
		markListCache.Clear()
	}()

	go func() {
		defer wg.Done()
		statisticsIndividualCache.Clear()
	}()

	go func() {
		defer wg.Done()
		statisticsAdminCache.Clear()
	}()
	wg.Wait()

	return updatedPost, nil
}

func (m *markListRepository) DeleteOneByUser(ctx context.Context, user primitive.ObjectID, markListID string) error {
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

	collectionMarkList := m.database.Collection(m.collectionMarkList)

	// Convert courseID string to ObjectID
	objID, err := primitive.ObjectIDFromHex(markListID)
	if err != nil {
		return err
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

	wg.Add(4)
	go func() {
		defer wg.Done()
		markListResCache.Clear()
	}()

	go func() {
		defer wg.Done()
		markListCache.Clear()
	}()

	go func() {
		defer wg.Done()
		statisticsIndividualCache.Clear()
	}()

	go func() {
		defer wg.Done()
		statisticsAdminCache.Clear()
	}()
	wg.Wait()

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

func (m *markListRepository) FetchManyByAdmin(ctx context.Context) (mark_list_domain.Response, error) {
	errCh := make(chan error, 1)
	markListCh := make(chan mark_list_domain.Response, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := markListResCache.Get("admin")
		if found {
			markListCh <- data
		}
	}()

	go func() {
		defer close(markListCh)
		wg.Wait()
	}()

	markListData := <-markListCh
	if !internal.IsZeroValue(markListData) {
		return markListData, nil
	}

	collectionMarkList := m.database.Collection(m.collectionMarkList)

	cursor, err := collectionMarkList.Find(ctx, bson.D{})
	if err != nil {
		return mark_list_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
		}
	}(cursor, ctx)

	var markLists []mark_list_domain.MarkList
	for cursor.Next(ctx) {
		var markList mark_list_domain.MarkList
		if err = cursor.Decode(&markList); err != nil {
			return mark_list_domain.Response{}, err
		}

		wg.Add(1)
		go func(markList mark_list_domain.MarkList) {
			defer wg.Done()
			// Thêm lesson vào slice lessons
			markLists = append(markLists, markList)
		}(markList)
	}
	wg.Wait()

	var statistics mark_list_domain.Statistics
	go func() {
		statistics, _ = m.Statistics(ctx)
	}()

	markListRes := mark_list_domain.Response{
		MarkList:   markLists,
		Statistics: statistics,
	}

	markListResCache.Set("admin", markListRes, cacheTTL)

	select {
	case err = <-errCh:
		return mark_list_domain.Response{}, err
	default:
		return markListRes, err
	}
}

func (m *markListRepository) Statistics(ctx context.Context) (mark_list_domain.Statistics, error) {
	statisticsAdminCh := make(chan mark_list_domain.Statistics, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := statisticsAdminCache.Get("admin")
		if found {
			statisticsAdminCh <- data
		}
	}()

	go func() {
		defer close(statisticsAdminCh)
		wg.Wait()
	}()

	statisticData := <-statisticsAdminCh
	if !internal.IsZeroValue(statisticData) {
		return statisticData, nil
	}

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

	statisticsAdminCache.Set("admin", statistics, cacheTTL)
	return statistics, nil
}

func (m *markListRepository) StatisticsIndividual(ctx context.Context, user primitive.ObjectID) (mark_list_domain.Statistics, error) {
	statisticIndividualCh := make(chan mark_list_domain.Statistics)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := statisticsIndividualCache.Get(user.Hex())
		if found {
			statisticIndividualCh <- data
		}
	}()

	go func() {
		defer close(statisticIndividualCh)
		wg.Wait()
	}()

	statisticData := <-statisticIndividualCh
	if !internal.IsZeroValue(statisticData) {
		return statisticData, nil
	}

	collectionMarkList := m.database.Collection(m.collectionMarkList)
	collectionMarkVocabulary := m.database.Collection(m.collectionMarkVocabulary)

	filter := bson.M{"user_id": user}
	countMarkList, err := collectionMarkList.CountDocuments(ctx, filter)
	if err != nil {
		return mark_list_domain.Statistics{}, err
	}

	countMarkVocabulary, err := collectionMarkVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return mark_list_domain.Statistics{}, err
	}

	statistics := mark_list_domain.Statistics{
		Total:           countMarkList,
		CountVocabulary: countMarkVocabulary,
	}

	statisticsIndividualCache.Set(user.Hex(), statistics, cacheTTL)
	return statistics, nil
}
