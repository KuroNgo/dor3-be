package exam_repository

import (
	exam_domain "clean-architecture/domain/exam"
	"clean-architecture/internal"
	"clean-architecture/internal/cache/memory"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type examRepository struct {
	database               *mongo.Database
	collectionLesson       string
	collectionUnit         string
	collectionExam         string
	collectionExamProcess  string
	collectionExamQuestion string
	collectionVocabulary   string
}

func NewExamRepository(db *mongo.Database, collectionExam string, collectionExamProcess string,
	collectionLesson string, collectionUnit string, collectionExamQuestion string, collectionVocabulary string) exam_domain.IExamRepository {
	return &examRepository{
		database:               db,
		collectionExam:         collectionExam,
		collectionExamProcess:  collectionExamProcess,
		collectionLesson:       collectionLesson,
		collectionUnit:         collectionUnit,
		collectionExamQuestion: collectionExamQuestion,
		collectionVocabulary:   collectionVocabulary,
	}
}

var (
	examCache        = memory.NewTTL[string, exam_domain.Exam]()
	examsCache       = memory.NewTTL[string, []exam_domain.Exam]()
	examProcessCache = memory.NewTTL[string, exam_domain.ExamProcessRes]()
	detailCache      = memory.NewTTL[string, exam_domain.DetailResponse]()
	statisticsCache  = memory.NewTTL[string, exam_domain.Statistics]()

	mu           sync.Mutex
	wg           sync.WaitGroup
	isProcessing bool
)

const (
	cacheTTL = 5 * time.Minute
)

func (e *examRepository) FetchOneByUnitIDInUser(ctx context.Context, userID primitive.ObjectID, unitID string) (exam_domain.ExamProcessRes, error) {
	errCh := make(chan error, 1)
	examProcessResCh := make(chan exam_domain.ExamProcessRes, 1)

	wg.Add(1)
	go func() {
		data, found := examProcessCache.Get(userID.Hex() + unitID)
		if found {
			examProcessResCh <- data
		}
	}()

	go func() {
		defer close(examProcessResCh)
		wg.Wait()
	}()

	examProcessResData := <-examProcessResCh
	if !internal.IsZeroValue(examProcessResData) {
		return examProcessResData, nil
	}

	collectionExam := e.database.Collection(e.collectionExam)
	collectionExamProcess := e.database.Collection(e.collectionExamProcess)

	idUnit, _ := primitive.ObjectIDFromHex(unitID)
	filter := bson.M{"unit_id": idUnit}
	filterExamProcess := bson.M{"user_id": userID, "unit_id": idUnit}

	countExam, err := collectionExam.CountDocuments(ctx, filter)
	if err != nil {
		return exam_domain.ExamProcessRes{}, err
	}

	count, err := collectionExamProcess.CountDocuments(ctx, filterExamProcess)
	if err != nil {
		return exam_domain.ExamProcessRes{}, err
	}

	if count < countExam {
		cursorExam, err := collectionExam.Find(ctx, filter)
		if err != nil {
			return exam_domain.ExamProcessRes{}, err
		}
		defer func(cursorExam *mongo.Cursor, ctx context.Context) {
			err = cursorExam.Close(ctx)
			if err != nil {
				errCh <- err
				return
			}
		}(cursorExam, ctx)

		for cursorExam.Next(ctx) {
			var exam exam_domain.Exam
			if err = cursorExam.Decode(&exam); err != nil {
				return exam_domain.ExamProcessRes{}, err
			}

			wg.Add(1)
			go func(exam exam_domain.Exam) {
				defer wg.Done()
				examProcess := exam_domain.ExamProcess{
					ExamID:     exam.ID,
					UserID:     userID,
					IsComplete: 0,
				}

				filterExam := bson.M{"exam_id": exam.ID, "user_id": userID}
				countExamChild, err := collectionExamProcess.CountDocuments(ctx, filterExam)
				if err != nil {
					errCh <- err
					return
				}

				if countExamChild == 0 {
					_, err = collectionExamProcess.InsertOne(ctx, &examProcess)
					if err != nil {
						errCh <- err
						return
					}
				}
			}(exam)
		}
		wg.Wait()
	}

	var exam exam_domain.Exam
	err = collectionExam.FindOne(ctx, filter).Decode(&exam)
	if err != nil {
		return exam_domain.ExamProcessRes{}, err
	}

	examProcessRes := exam_domain.ExamProcessRes{
		Exam:       exam,
		UserID:     userID,
		IsComplete: 0,
	}

	examProcessCache.Set(userID.Hex(), examProcessRes, cacheTTL)

	select {
	case err = <-errCh:
		return exam_domain.ExamProcessRes{}, err
	default:
		return examProcessRes, nil
	}
}

func (e *examRepository) UpdateCompletedInUser(ctx context.Context, userID primitive.ObjectID, exam *exam_domain.ExamProcess) error {
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

	collectionExamProcess := e.database.Collection(e.collectionExamProcess)

	filter := bson.D{{"_id", exam.ExamID}, {"user_id", userID}}
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := collectionExamProcess.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e *examRepository) FetchManyInAdmin(ctx context.Context, page string) ([]exam_domain.Exam, exam_domain.DetailResponse, error) {
	errCh := make(chan error, 1)
	examsCh := make(chan []exam_domain.Exam, 1)
	detailCh := make(chan exam_domain.DetailResponse, 1)

	wg.Add(2)
	go func() {
		defer wg.Done()
		data, found := examsCache.Get(page)
		if found {
			examsCh <- data
		}
	}()

	go func() {
		defer wg.Done()
		data, found := detailCache.Get("detail" + page)
		if found {
			detailCh <- data
		}
	}()

	go func() {
		defer close(examsCh)
		defer close(detailCh)
		wg.Wait()
	}()

	examsData := <-examsCh
	detailData := <-detailCh
	if !internal.IsZeroValue(examsData) && !internal.IsZeroValue(detailData) {
		return examsData, detailData, nil
	}

	collectionExam := e.database.Collection(e.collectionExam)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))
	count, err := collectionExam.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	// Tính toán tổng số trang dựa trên số lượng khóa học và số khóa học mỗi trang
	totalPages := (count + int64(perPage) - 1) / int64(perPage)

	cursor, err := collectionExam.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var exams []exam_domain.Exam
	exams = make([]exam_domain.Exam, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var exam exam_domain.Exam
		if err = cursor.Decode(&exam); err != nil {
			return nil, exam_domain.DetailResponse{}, err
		}

		wg.Add(1)
		go func(exam exam_domain.Exam) {
			defer wg.Done()
			exams = append(exams, exam)
		}(exam)
	}
	wg.Wait()

	statisticsCh := make(chan exam_domain.Statistics)
	go func() {
		statistics, _ := e.Statistics(ctx)
		statisticsCh <- statistics
	}()
	statistics := <-statisticsCh

	detail := exam_domain.DetailResponse{
		Page:        totalPages,
		CurrentPage: pageNumber,
		CountExam:   int64(len(exams)),
		Statistics:  statistics,
	}

	examsCache.Set(page, exams, cacheTTL)
	detailCache.Set("detail"+page, detail, cacheTTL)

	select {
	case err = <-errCh:
		return nil, exam_domain.DetailResponse{}, err
	default:
		return exams, detail, nil
	}
}

func (e *examRepository) FetchManyByUnitIDInAdmin(ctx context.Context, unitID string, page string) ([]exam_domain.Exam, exam_domain.DetailResponse, error) {
	errCh := make(chan error, 1)
	examsCh := make(chan []exam_domain.Exam, 1)
	detailCh := make(chan exam_domain.DetailResponse, 1)

	wg.Add(2)
	go func() {
		defer wg.Done()
		data, found := examsCache.Get(unitID + page)
		if found {
			examsCh <- data
		}
	}()

	go func() {
		defer wg.Done()
		data, found := detailCache.Get("detail" + unitID)
		if found {
			detailCh <- data
		}
	}()

	go func() {
		defer close(examsCh)
		defer close(detailCh)
		wg.Wait()
	}()

	examsData := <-examsCh
	detailData := <-detailCh
	if !internal.IsZeroValue(examsData) && !internal.IsZeroValue(detailData) {
		return examsData, detailData, nil
	}

	collectionExam := e.database.Collection(e.collectionExam)

	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		return nil, exam_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip)).SetSort(bson.D{{"_id", -1}})

	// Convert unitID to ObjectID
	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	// Count documents for pagination
	count, err := collectionExam.CountDocuments(ctx, bson.M{"unit_id": idUnit})
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}
	totalPages := (count + int64(perPage) - 1) / int64(perPage) // Calculate total pages

	// Query for exercises
	filter := bson.M{"unit_id": idUnit}
	cursor, err := collectionExam.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var exams []exam_domain.Exam
	exams = make([]exam_domain.Exam, 0, cursor.RemainingBatchLength())
	// Process each exercise
	for cursor.Next(ctx) {
		var exam exam_domain.Exam
		if err = cursor.Decode(&exam); err != nil {
			return nil, exam_domain.DetailResponse{}, err
		}

		wg.Add(1)
		go func(exam exam_domain.Exam) {
			defer wg.Done()
			exams = append(exams, exam)
		}(exam)
	}
	wg.Wait()

	if err = cursor.Err(); err != nil {
		return nil, exam_domain.DetailResponse{}, err
	}

	statisticsCh := make(chan exam_domain.Statistics)
	go func() {
		statistics, _ := e.Statistics(ctx)
		statisticsCh <- statistics
	}()
	statistics := <-statisticsCh

	detail := exam_domain.DetailResponse{
		CountExam:   int64(len(exams)),
		Page:        totalPages,
		Statistics:  statistics,
		CurrentPage: pageNumber,
	}

	examsCache.Set(unitID+page, exams, cacheTTL)
	detailCache.Set("detail"+unitID, detail, cacheTTL)

	select {
	case err = <-errCh:
		return nil, exam_domain.DetailResponse{}, err
	default:
		return exams, detail, nil
	}
}

func (e *examRepository) FetchExamByIDInAdmin(ctx context.Context, id string) (exam_domain.Exam, error) {
	examCh := make(chan exam_domain.Exam, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := examCache.Get(id)
		if found {
			examCh <- data
		}
	}()

	go func() {
		defer close(examCh)
		wg.Wait()
	}()

	examData := <-examCh
	if !internal.IsZeroValue(examData) {
		return examData, nil
	}

	collectionExam := e.database.Collection(e.collectionExam)

	idExam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return exam_domain.Exam{}, err
	}

	var exam exam_domain.Exam
	filter := bson.M{"_id": idExam}
	err = collectionExam.FindOne(ctx, filter).Decode(&exam)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return exam_domain.Exam{}, errors.New("exam not found")
		}
		return exam_domain.Exam{}, err
	}

	examCache.Set(id, exam, cacheTTL)
	return exam, nil
}

func (e *examRepository) FetchOneByUnitIDInAdmin(ctx context.Context, unitID string) (exam_domain.Exam, error) {
	errCh := make(chan error, 1)
	examCh := make(chan exam_domain.Exam, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := examCache.Get(unitID)
		if found {
			examCh <- data
		}
	}()

	go func() {
		defer close(examCh)
		wg.Wait()
	}()

	examData := <-examCh
	if !internal.IsZeroValue(examData) {
		return examData, nil
	}

	collectionExam := e.database.Collection(e.collectionExam)

	idUnit, err := primitive.ObjectIDFromHex(unitID)
	if err != nil {
		return exam_domain.Exam{}, err
	}

	filter := bson.M{"unit_id": idUnit}
	cursor, err := collectionExam.Find(ctx, filter)
	if err != nil {
		return exam_domain.Exam{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var exams []exam_domain.Exam
	for cursor.Next(ctx) {
		var exam exam_domain.Exam
		if err = cursor.Decode(&exam); err != nil {
			return exam_domain.Exam{}, err
		}

		wg.Add(1)
		go func(exam exam_domain.Exam) {
			defer wg.Done()
			exams = append(exams, exam)
		}(exam)
	}
	wg.Wait()

	// Kiểm tra nếu danh sách exercises không rỗng
	if len(exams) == 0 {
		return exam_domain.Exam{}, errors.New("no exercises found")
	}

	// Chọn một giá trị ngẫu nhiên từ danh sách exercises
	randomIndex := rand.Intn(len(exams))
	randomExam := exams[randomIndex]

	examCache.Set(unitID, randomExam, cacheTTL)

	select {
	case err = <-errCh:
		return exam_domain.Exam{}, err
	default:
		return randomExam, nil
	}
}

func (e *examRepository) CreateOneInAdmin(ctx context.Context, exam *exam_domain.Exam) error {
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

	collectionExam := e.database.Collection(e.collectionExam)

	collectionLesson := e.database.Collection(e.collectionLesson)
	collectionUnit := e.database.Collection(e.collectionUnit)

	filterLessonID := bson.M{"_id": exam.LessonID}
	countLessonID, err := collectionLesson.CountDocuments(ctx, filterLessonID)
	if err != nil {
		return err
	}
	if countLessonID == 0 {
		return errors.New("the lesson ID does not exist")
	}

	filterUnitID := bson.M{"_id": exam.UnitID, "lesson_id": exam.LessonID}
	countUnitID, err := collectionUnit.CountDocuments(ctx, filterUnitID)
	if err != nil {
		return err
	}
	if countUnitID == 0 {
		return errors.New("the unit ID does not exist")
	}

	_, err = collectionExam.InsertOne(ctx, exam)
	if err != nil {
		return err
	}

	return nil
}

func (e *examRepository) UpdateOneInAdmin(ctx context.Context, exam *exam_domain.Exam) (*mongo.UpdateResult, error) {
	collectionExam := e.database.Collection(e.collectionExam)
	collectionExamProcess := e.database.Collection(e.collectionExamProcess)

	filter := bson.M{"_id": exam.ID}
	update := bson.M{
		"$set": bson.M{
			"description": exam.Description,
			"title":       exam.Title,
			"duration":    exam.Duration,
		},
	}

	data, err := collectionExam.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	_, err = collectionExamProcess.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (e *examRepository) DeleteOneInAdmin(ctx context.Context, examID string) error {
	collectionExam := e.database.Collection(e.collectionExam)
	collectionExamProcess := e.database.Collection(e.collectionExamProcess)

	objID, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionExam.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exam is removed or have not exist`)
	}

	_, err = collectionExam.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	_, err = collectionExamProcess.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return err
}

func (e *examRepository) countQuestion(ctx context.Context, examID string) int64 {
	collectionExamQuestion := e.database.Collection(e.collectionExamQuestion)

	idExam, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return 0
	}

	filter := bson.M{"exam_id": idExam}
	count, err := collectionExamQuestion.CountDocuments(ctx, filter)
	if err != nil {
		return 0
	}

	return count
}

func (e *examRepository) Statistics(ctx context.Context) (exam_domain.Statistics, error) {
	statisticsCh := make(chan exam_domain.Statistics)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := statisticsCache.Get("statistics")
		if found {
			statisticsCh <- data
		}
	}()

	go func() {
		defer close(statisticsCh)
		wg.Wait()
	}()

	collectionExam := e.database.Collection(e.collectionExam)

	count, err := collectionExam.CountDocuments(ctx, bson.D{})
	if err != nil {
		return exam_domain.Statistics{}, err
	}

	statistics := exam_domain.Statistics{
		Total: count,
	}

	statisticsCache.Set("statistics", statistics, cacheTTL)
	return statistics, nil
}
