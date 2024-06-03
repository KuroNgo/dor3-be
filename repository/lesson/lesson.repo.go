package lesson_repository

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
	"strconv"
	"sync"
	"time"
)

type lessonRepository struct {
	database             *mongo.Database
	collectionLesson     string
	collectionCourse     string
	collectionUnit       string
	collectionVocabulary string
}

// NewLessonRepository hàm khởi tạo (constructor) để khởi tạo instance của struct
func NewLessonRepository(db *mongo.Database, collectionLesson string, collectionCourse string, collectionUnit string, collectionVocabulary string) lesson_domain.ILessonRepository {
	return &lessonRepository{
		database:             db,
		collectionLesson:     collectionLesson,
		collectionCourse:     collectionCourse,
		collectionUnit:       collectionUnit,
		collectionVocabulary: collectionVocabulary,
	}
}

var (
	lessonsCache       = cache.NewTTL[string, []lesson_domain.LessonResponse]()
	lessonCache        = cache.NewTTL[string, lesson_domain.LessonResponse]()
	detailCache        = cache.NewTTL[string, lesson_domain.DetailResponse]()
	coursePrimOIDCache = cache.NewTTL[string, primitive.ObjectID]()
	statisticsCache    = cache.NewTTL[string, lesson_domain.Statistics]()

	wg sync.WaitGroup
	mu sync.Mutex
)

// FetchMany lấy tất cả bài học (lesson) cùng một lúc (concurrency).
// Hàm này nhận vào số trang (page) và trả về một mảng bài học làm khóa và nội dung của bài học tương ứng làm giá trị.
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả với các kết quả đã lấy được
// FIXME: thực hiện gắn lỗi vào channel giúp tối ưu hóa xử lý
func (l *lessonRepository) FetchMany(ctx context.Context, page string) ([]lesson_domain.LessonResponse, lesson_domain.DetailResponse, error) {
	// Khởi tạo channel để luu trữ lỗi
	errCh := make(chan error)
	// Khởi tạo channel để lưu trữ kết quả lesson
	lessonsCh := make(chan []lesson_domain.LessonResponse, 1)
	// Khởi tạo channel để lưu trữ kết quả detail
	detailCh := make(chan lesson_domain.DetailResponse, 1)
	// Sử dụng WaitGroup để đợi tất cả các goroutine hoàn thành
	wg.Add(2)
	// Khởi động một goroutine cho tìm dữ liệu lesson trong cache memory
	go func() {
		defer wg.Done()
		data, found := lessonsCache.Get(page)
		if found {
			lessonsCh <- data
			return
		}
	}()

	// Khởi động một goroutine cho tìm dữ liệu detail trong cache memory
	go func() {
		defer wg.Done()
		data, found := detailCache.Get("detail")
		if found {
			detailCh <- data
			return
		}
	}()

	// Goroutine để đóng các channel khi tất cả các công việc hoàn thành
	go func() {
		defer close(lessonsCh)
		defer close(detailCh)
		wg.Wait()
	}()

	// Gán giá trị từ channel
	lessonData := <-lessonsCh
	detailData := <-detailCh

	// kiểm tra dữ liệu Data có rỗng hay không,
	// nếu không sẽ trả về dữ lệu trong cache vừa tìm được
	// Ngược lại, sẽ thực hiện quy trình tìm
	if !internal.IsZeroValue(lessonData) && !internal.IsZeroValue(detailData) {
		return lessonData, detailData, nil
	}

	// khởi tạo đối tượng collection, ở đây là lesson và unit (lesson tham chiếu)
	collectionLesson := l.database.Collection(l.collectionLesson)
	collectionUnit := l.database.Collection(l.collectionUnit)

	// thực hiện phân trang
	// lấy số trang từ client
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, errors.New("invalid page number")
	}

	// tối đa dữ liệu gửi đến ở mỗi trang được yêu cầu
	perPage := 10
	// nếu dữ liệu nhỏ hơn 1 sẽ skip
	skip := (pageNumber - 1) * perPage
	// thực hiện các yêu cầu đã neu ở trên
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	// đếm số lượng lesson có trong project
	count, err := collectionLesson.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}

	// đếm số lượng trang
	totalPages := (count + int64(perPage) - 1) / int64(perPage)

	// thực hiện tìm kiếm theo điêu kiện options
	cursor, err := collectionLesson.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}
	// thực hiện close cursor nếu hàm đã hoàn thành hoặc bị lỗi
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var lessons []lesson_domain.LessonResponse
	wg.Add(1)
	// Khởi động một goroutine cho mỗi cursor
	// TODO: Xử lý tìm từng dữ liệu liên quan đến lesson bao gồm các thông tin cơ bản và các thông tin khác (thống kê, số lượng unit đã hoàn thành của user)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var arrIsComplete []int
			// chuyển đổi sang JSON cho lesson
			var lesson lesson_domain.LessonResponse
			if err = cursor.Decode(&lesson); err != nil {
				errCh <- err
				return
			}

			// từ lesson, tìm những unit thuộc lesson
			var units []unit_domain.UnitResponse
			filterLesson := bson.M{"lesson_id": lesson.ID}
			cursorUnit, err := collectionUnit.Find(ctx, filterLesson)
			if err != nil {
				errCh <- err
				return
			}

			for cursorUnit.Next(ctx) {
				var unit unit_domain.UnitResponse
				err := cursorUnit.Decode(&unit)
				if err != nil {
					errCh <- err
					return
				}

				units = append(units, unit)
			}

			// từ unit, lấy gia tri complete
			for _, unit := range units {
				arrIsComplete = append(arrIsComplete, unit.IsComplete)
			}

			// gắn giá trị complete của từng unit vào mảng và gưi giá trị đó cho lesson
			lesson.UnitIsComplete = arrIsComplete
			// tạo channel để thu thập kết quả
			countVocabularyCh := make(chan int32)
			countUnitCh := make(chan int32)

			// Goroutine giúp lấy kết quả đếm unit cho chủ đề
			go func() {
				defer close(countUnitCh)
				// Lấy thông tin liên quan cho mỗi chủ đề
				countUnit, err := l.countUnitsByLessonsID(ctx, lesson.ID)
				if err != nil {
					errCh <- err
					return
				}

				// gắn giá trị đếm vào channel
				countUnitCh <- countUnit
			}()

			// Goroutine giúp lấy kết quả đếm từ vựng cho chủ đề
			go func() {
				defer close(countVocabularyCh)
				countVocabulary, err := l.countVocabularyByLessonID(ctx, lesson.ID)
				if err != nil {
					errCh <- err
					return
				}

				// gắn giá trị đếm vào channel
				countVocabularyCh <- countVocabulary
			}()

			countUnit := <-countUnitCh
			countVocabulary := <-countVocabularyCh

			lesson.CountUnit = countUnit
			lesson.CountVocabulary = countVocabulary

			// Thêm lesson vào slice lessons
			lessons = append(lessons, lesson)
		}
	}()
	wg.Wait()

	// Channel để thu thập kết quả thống kê
	statisticsCh := make(chan lesson_domain.Statistics)
	// Goroutine thực hiện lấy giá trị thống kê toàn bộ
	go func() {
		defer close(statisticsCh)
		statistics, _ := l.Statistics(ctx)
		statisticsCh <- statistics
	}()
	statistics := <-statisticsCh

	// Luu tất cả thông tin liên quan vào detail
	response := lesson_domain.DetailResponse{
		Page:        totalPages,
		CurrentPage: pageNumber,
		Statistics:  statistics,
	}

	// Thiết lập Set cache memory với dữ liệu cần thiết với thơi gian là 5 phút
	lessonsCache.Set(page, lessons, 5*time.Minute)
	detailCache.Set("detail", response, 5*time.Minute)

	// Thu thập kết quả
	select {
	// Nếu có lỗi, sẽ thực hiện trả về lỗi
	case err = <-errCh:
		return nil, lesson_domain.DetailResponse{}, err
	// Ngược lại, sẽ trả về giá trị
	default:
		return lessons, response, err
	}
}

// FetchManyNotPagination lấy tất cả bài học (lesson) cùng một lúc (concurrency)
// Hàm này không nhận đầu vào (input) và trả về một mảng bài học làm khóa và nội dung của bài học tương ứng làm giá trị
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả về với các kết quả đã lấy được
func (l *lessonRepository) FetchManyNotPagination(ctx context.Context) ([]lesson_domain.LessonResponse, lesson_domain.DetailResponse, error) {
	// Khởi tạo channels để lưu trữ lỗi, bài học và chi tiết
	errCh := make(chan error)
	lessonsCh := make(chan []lesson_domain.LessonResponse, 1)
	detailCh := make(chan lesson_domain.DetailResponse, 1)
	wg.Add(2)

	// Goroutine để lấy bài học từ cache
	go func() {
		defer wg.Done()
		data, found := lessonsCache.Get("lessons")
		if found {
			lessonsCh <- data
			return
		}
	}()

	// Goroutine để lấy chi tiết từ cache
	go func() {
		defer wg.Done()
		data, found := detailCache.Get("detail")
		if found {
			detailCh <- data
			return
		}
	}()

	// Goroutine để đóng channels sau khi tất cả các công việc hoàn thành
	go func() {
		defer close(lessonsCh)
		defer close(detailCh)
		wg.Wait()
	}()

	// Lấy bài học và chi tiết từ các channels
	lessonsData := <-lessonsCh
	detailData := <-detailCh

	// Kiểm tra xem dữ liệu có tồn tại trong cache không, nếu không thì lấy từ database
	if !internal.IsZeroValue(lessonsData) && !internal.IsZeroValue(detailData) {
		return lessonsData, detailData, nil
	}

	// Lấy các collection của bài học và chủ đề
	collectionLesson := l.database.Collection(l.collectionLesson)
	collectionUnit := l.database.Collection(l.collectionUnit)

	// Lấy các bài học từ database
	cursor, err := collectionLesson.Find(ctx, bson.D{})
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	// Khởi tạo slice để lưu trữ bài học
	var lessons []lesson_domain.LessonResponse
	wg.Add(1)

	// Goroutine để lấy bài học và các đơn vị liên quan của chúng đồng thời
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var arrIsComplete []int
			var lesson lesson_domain.LessonResponse
			if err = cursor.Decode(&lesson); err != nil {
				errCh <- err
				return
			}

			// Lấy các đơn vị liên quan đến bài học
			var units []unit_domain.UnitResponse
			filterLesson := bson.M{"lesson_id": lesson.ID}
			cursorUnit, err := collectionUnit.Find(ctx, filterLesson)
			if err != nil {
				errCh <- err
				return
			}

			// Lặp qua các đơn vị và thêm trạng thái hoàn thành vào một mảng
			for cursorUnit.Next(ctx) {
				var unit unit_domain.UnitResponse
				err = cursorUnit.Decode(&unit)
				if err != nil {
					errCh <- err
					return
				}
				units = append(units, unit)
			}

			for _, unit := range units {
				arrIsComplete = append(arrIsComplete, unit.IsComplete)
			}

			lesson.UnitIsComplete = arrIsComplete

			// Lấy thông tin bổ sung cho mỗi bài học
			countUnitCh := make(chan int32)
			go func() {
				defer close(countUnitCh)
				countUnit, err := l.countUnitsByLessonsID(ctx, lesson.ID)
				if err != nil {
					errCh <- err
					return
				}
				countUnitCh <- countUnit
			}()

			countVocabularyCh := make(chan int32)
			go func() {
				defer close(countVocabularyCh)
				countVocabulary, err := l.countVocabularyByLessonID(ctx, lesson.ID)
				if err != nil {
					errCh <- err
					return
				}
				countVocabularyCh <- countVocabulary
			}()

			countUnit := <-countUnitCh
			countVocabulary := <-countVocabularyCh

			lesson.CountUnit = countUnit
			lesson.CountVocabulary = countVocabulary

			lessons = append(lessons, lesson)
		}
	}()

	wg.Wait()

	// Lấy thống kê
	var statisticsCh = make(chan lesson_domain.Statistics)
	go func() {
		defer close(statisticsCh)
		statistic, _ := l.Statistics(ctx)
		statisticsCh <- statistic
	}()
	statisticsData := <-statisticsCh

	// Kết hợp dữ liệu đã lấy thành phản hồi chi tiết
	detail := lesson_domain.DetailResponse{
		Statistics: statisticsData,
	}

	// Lưu trữ dữ liệu đã lấy vào cache
	lessonsCache.Set("lessons", lessons, 5*time.Minute)
	detailCache.Set("detail", detail, 5*time.Minute)

	// Trả về dữ liệu đã lấy hoặc lỗi
	select {
	case err = <-errCh:
		return nil, lesson_domain.DetailResponse{}, err
	default:
		return lessons, detail, nil
	}
}

// FetchByID lấy bài học (lesson) theo ID
// Hàm này nhận đầu vào (input) là lessonID và trả về một bài học làm khóa và nội dung cuủa bài học tương ứng làm giá trị
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả về với các kết quả đã lấy được
func (l *lessonRepository) FetchByID(ctx context.Context, lessonID string) (lesson_domain.LessonResponse, error) {
	// Khởi tạo channel để luu trữ lỗi
	errCh := make(chan error)
	// Khởi tạo channel để lưu trữ kết quả lesson
	lessonCh := make(chan lesson_domain.LessonResponse)
	// Sử dụng waitGroup để đợi tất cả goroutine hoàn thành
	wg.Add(1)
	// Khởi động Goroutine giúp tìm dữ liệu lesson
	// theo id trong cache (đã từng tìm lessonID này hay chưa)
	go func() {
		defer wg.Done()
		data, found := lessonCache.Get(lessonID)
		if found {
			lessonCh <- data
		}
	}()

	// Goroutine để đóng các channel khi tất cả các công việc hoàn thành
	go func() {
		defer close(lessonCh)
		wg.Wait()

	}()

	// Channel gửi giá trị cho biến lessonData
	lessonData := <-lessonCh
	// Kiểm tra giá trị lessonData có null ?
	// Nếu không thì sẽ thực hiện trả về giá trị
	// Ngược lại thì thực hiện tìm theo LessonID
	if !internal.IsZeroValue(lessonData) {
		return lessonData, nil
	}

	collectionLesson := l.database.Collection(l.collectionLesson)
	collectionUnit := l.database.Collection(l.collectionUnit)

	// Thực hiện chuyển đổi lessonID từ string sang primitive.ObjectID
	idLesson, err := primitive.ObjectIDFromHex(lessonID)
	if err != nil {
		return lesson_domain.LessonResponse{}, err
	}

	// Lấy dữ liệu lessonID vừa chuyển đổi, thực hiện tìm kiếm theo id
	filter := bson.M{"_id": idLesson}

	var lesson lesson_domain.LessonResponse
	// Thực hiện tìm kiếm lesson theo id
	err = collectionLesson.FindOne(ctx, filter).Decode(&lesson)
	if err != nil {
		return lesson_domain.LessonResponse{}, err
	}

	// từ lesson, tìm những unit thuộc lesson
	var units []unit_domain.UnitResponse
	filterLesson := bson.M{"lesson_id": lesson.ID}
	cursorUnit, err := collectionUnit.Find(ctx, filterLesson)
	if err != nil {
		return lesson_domain.LessonResponse{}, err
	}

	for cursorUnit.Next(ctx) {
		var unit unit_domain.UnitResponse
		err = cursorUnit.Decode(&unit)
		if err != nil {
			return lesson_domain.LessonResponse{}, err
		}

		units = append(units, unit)
	}

	var arrIsComplete []int
	// từ unit, lấy gia tri complete
	for _, unit := range units {
		arrIsComplete = append(arrIsComplete, unit.IsComplete)
	}

	// gắn giá trị complete của từng unit vào mảng và gưi giá trị đó cho lesson
	lesson.UnitIsComplete = arrIsComplete

	countUnitCh := make(chan int32)
	countVocabularyCh := make(chan int32)

	// Goroutine để thực hiên đếm số lượng vocabulary trong lesson
	go func() {
		defer close(countVocabularyCh)
		countVocabulary, err := l.countVocabularyByLessonID(ctx, lesson.ID)
		if err != nil {
			errCh <- err
			return
		}
		countVocabularyCh <- countVocabulary
	}()

	// Goroutine để thực hiên đếm số lượng unit trong lesson
	go func() {
		defer close(countUnitCh)
		countUnit, err := l.countUnitsByLessonsID(ctx, lesson.ID)
		if err != nil {
			errCh <- err
			return
		}
		countUnitCh <- countUnit
	}()

	// Channel gửi giá trị, sau đó lesson sẽ nhận giá trị tương ứng
	countUnit := <-countUnitCh
	countVocabulary := <-countVocabularyCh
	lesson.CountVocabulary = countVocabulary
	lesson.CountUnit = countUnit

	// Thiết lập Set cache memory với dữ liệu cần thiết với thơi gian là 5 phút
	lessonCache.Set(lessonID, lesson, 5*time.Minute)

	// Thu thập kết quả
	select {
	// Nếu có lỗi, sẽ thực hiện trả về lỗi
	case err = <-errCh:
		return lesson_domain.LessonResponse{}, err
	// Ngược lại, sẽ trả về giá trị
	default:
		return lesson, nil
	}
}

// FetchByIdCourse lấy bài học (lesson) theo courseID
// Hàm này nhận tham số là idCourse và page và trả v một mảng bài học làm khóa và thống kê đi kèm (detail)
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả về với các kết quả đã lấy được
func (l *lessonRepository) FetchByIdCourse(ctx context.Context, idCourse string, page string) ([]lesson_domain.LessonResponse, lesson_domain.DetailResponse, error) {
	errCh := make(chan error)
	lessonsCh := make(chan []lesson_domain.LessonResponse, 1)
	detailCh := make(chan lesson_domain.DetailResponse, 1)

	wg.Add(2)
	go func() {
		defer wg.Done()
		data, found := lessonsCache.Get(idCourse + page)
		if found {
			lessonsCh <- data
			return
		}
	}()

	go func() {
		defer wg.Done()
		data, found := detailCache.Get(idCourse + "detail")
		if found {
			detailCh <- data
			return
		}
	}()

	go func() {
		defer close(lessonsCh)
		defer close(detailCh)
		wg.Wait()
	}()

	lessonData := <-lessonsCh
	detailData := <-detailCh
	if !internal.IsZeroValue(lessonData) && !internal.IsZeroValue(detailData) {
		return lessonData, detailData, nil
	}

	collectionLesson := l.database.Collection(l.collectionLesson)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, errors.New("invalid page number")
	}
	perPage := 10
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	calCh := make(chan int64)
	go func() {
		count, err := collectionLesson.CountDocuments(ctx, bson.D{})
		if err != nil {
			errCh <- err
			return
		}

		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calCh <- cal1 + 1
		}
	}()

	idCourse2, err := primitive.ObjectIDFromHex(idCourse)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}

	filter := bson.M{"course_id": idCourse2}

	cursor, err := collectionLesson.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, lesson_domain.DetailResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var lessons []lesson_domain.LessonResponse
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			var lesson lesson_domain.LessonResponse
			if err = cursor.Decode(&lesson); err != nil {
				errCh <- err
				return
			}

			countUnitCh := make(chan int32)
			go func() {
				defer close(countUnitCh)
				// Lấy thông tin liên quan cho mỗi chủ đề
				countUnit, err := l.countUnitsByLessonsID(ctx, lesson.ID)
				if err != nil {
					errCh <- err
					return
				}

				countUnitCh <- countUnit
			}()

			countVocabularyCh := make(chan int32)
			go func() {
				defer close(countVocabularyCh)
				countVocabulary, err := l.countVocabularyByLessonID(ctx, lesson.ID)
				if err != nil {
					errCh <- err
					return
				}

				countVocabularyCh <- countVocabulary
			}()

			countUnit := <-countUnitCh

			countVocabulary := <-countVocabularyCh

			// Gắn CourseID vào bài học
			lesson.CourseID = idCourse2
			lesson.CountVocabulary = countVocabulary
			lesson.CountUnit = countUnit

			lessons = append(lessons, lesson)
		}
	}()
	wg.Wait()

	cal := <-calCh
	response := lesson_domain.DetailResponse{
		Page:        cal,
		CurrentPage: pageNumber,
	}

	lessonsCache.Set(idCourse+page, lessons, 5*time.Minute)
	detailCache.Set(idCourse+"detail", response, 5*time.Minute)

	select {
	case err = <-errCh:
		return nil, lesson_domain.DetailResponse{}, err
	default:
		return lessons, response, nil
	}
}

// FindCourseIDByCourseName lấy khóa học lấy mã khóa học (courseid) theo courseName
// Hàm này nhận tham số là courseNam và trả về một oid (primitive.ObjectID)
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả về với kết quả lấy được
// TODO: Hàm này chỉ dùng để hỗ trợ (helper) không làm api (controller)
func (l *lessonRepository) FindCourseIDByCourseName(ctx context.Context, courseName string) (primitive.ObjectID, error) {
	courseIDCh := make(chan primitive.ObjectID)
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := coursePrimOIDCache.Get(courseName)
		if found {
			courseIDCh <- data
			return
		}
	}()

	go func() {
		defer close(courseIDCh)
		wg.Wait()
	}()

	courseIDData := <-courseIDCh
	if !internal.IsZeroValue(courseIDData) {
		return courseIDData, nil
	}

	collectionCourse := l.database.Collection(l.collectionCourse)

	filter := bson.M{"name": courseName}
	var data struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	err := collectionCourse.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return primitive.NilObjectID, err
	}

	coursePrimOIDCache.Set(courseName, data.Id, 5*time.Minute)
	return data.Id, nil
}

// CreateOne khởi tạo bài học (lesson) theo đối tượng lesson
// Hàm này nhận tham số là một đối tượng và trả về kết quả thông tin xử lý (không phải là thông tin của đối tượng đó)
// Nếu có lỗi xảy ra trong quá trình xử lý, lỗi sẽ được trả về với kết quả đã lấy được và dừng chương trình
func (l *lessonRepository) CreateOne(ctx context.Context, lesson *lesson_domain.Lesson) error {
	collectionLesson := l.database.Collection(l.collectionLesson)
	collectionCourse := l.database.Collection(l.collectionCourse)

	filter := bson.M{"name": lesson.Name}

	// check exists with CountDocuments
	count, err := collectionLesson.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the lesson name did exist")
	}

	filterCourse := bson.M{"_id": lesson.CourseID}
	countParent, err := collectionCourse.CountDocuments(ctx, filterCourse)
	if err != nil {
		return err
	}
	if countParent == 0 {
		return errors.New("the course ID do not exist")
	}

	_, err = collectionLesson.InsertOne(ctx, lesson)

	// Clear data value in cache memory
	wg.Add(2)
	go func() {
		defer wg.Done()
		lessonsCache.Clear()
	}()

	// clear data value in cache memory due to increase num
	go func() {
		defer wg.Done()
		detailCache.Clear()
	}()

	wg.Wait()
	return nil
}

func (l *lessonRepository) CreateOneByNameCourse(ctx context.Context, lesson *lesson_domain.Lesson) error {
	collectionLesson := l.database.Collection(l.collectionLesson)

	filter := bson.M{"name": lesson.Name}
	count, err := collectionLesson.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the lesson name did exist")
	}

	_, err = collectionLesson.InsertOne(ctx, lesson)
	// Clear data value in cache memory
	wg.Add(2)
	go func() {
		defer wg.Done()
		lessonsCache.Clear()
	}()

	// clear data value in cache memory due to increase num
	go func() {
		defer wg.Done()
		detailCache.Clear()
	}()

	wg.Wait()
	return nil
}

// UpdateOne cập nhật bài học (lesson) theo đối tượng lesson
// Hàm nhận tham số là đối tượng lesson. Nếu thành công sẽ trả về thông tin cập nhật (không phải thông tin đối tượng).
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả về với kết quả đã lấy được
func (l *lessonRepository) UpdateOne(ctx context.Context, lesson *lesson_domain.Lesson) (*mongo.UpdateResult, error) {
	collection := l.database.Collection(l.collectionLesson)

	filter := bson.M{"_id": lesson.ID}

	update := bson.M{
		"$set": bson.M{
			"name":        lesson.Name,
			"content":     lesson.Content,
			"image_url":   lesson.ImageURL,
			"updated_at":  lesson.UpdatedAt,
			"who_updates": lesson.WhoUpdates,
		},
	}

	mu.Lock()
	data, err := collection.UpdateOne(ctx, filter, &update)
	mu.Unlock()
	if err != nil {
		return nil, err
	}

	// Clear data value in cache memory for courses
	wg.Add(2)
	go func() {
		defer wg.Done()
		lessonsCache.Clear()
	}()

	// clear data value with id courseID in cache
	go func() {
		defer wg.Done()
		lessonCache.Remove(lesson.ID.Hex())
	}()
	wg.Wait()

	return data, err
}

// UpdateComplete cập nhật mức độ hoàn thành bài học (lesson) theo đối tượng course
// Hàm nhận tham số là đối tượng lesson. Nếu thành công sẽ trả về thông tin cập nhật (không phải thông tin đối tượng).
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả về với kết quả đã lấy được
// Hàm có sử dụng concurrency, giúp xử lý các tác vụ về người dùng hiệu quả
func (l *lessonRepository) UpdateComplete(ctx context.Context, lessonID string, lesson lesson_domain.Lesson) error {
	collection := l.database.Collection(l.collectionUnit)

	filter := bson.D{{Key: "_id", Value: lessonID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "is_complete", Value: lesson.IsCompleted},
		{Key: "who_updates", Value: lesson.WhoUpdates},
	}}}

	_, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return err
	}

	// Clear data value in cache memory for courses
	wg.Add(2)
	go func() {
		defer wg.Done()
		lessonsCache.Clear()
	}()

	// clear data value with id courseID in cache
	go func() {
		defer wg.Done()
		lessonCache.Remove(lesson.ID.Hex())
	}()
	wg.Wait()

	return nil
}

// UpdateImage cập nhật bài học (lesson) theo đối tượng lesson
// Hàm nhận tham số là file image. Nếu thành công sẽ trả về thông tin cập nhật (không phải thông tin đối tượng).
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả về với kết quả đã lấy được
func (l *lessonRepository) UpdateImage(ctx context.Context, lesson *lesson_domain.Lesson) (*mongo.UpdateResult, error) {
	collection := l.database.Collection(l.collectionLesson)

	filter := bson.M{"_id": lesson.ID}

	update := bson.M{
		"$set": bson.M{
			"image_url":   lesson.ImageURL,
			"asset_url":   lesson.AssetURL,
			"updated_at":  lesson.UpdatedAt,
			"who_updates": lesson.WhoUpdates,
		},
	}

	mu.Lock()
	data, err := collection.UpdateOne(ctx, filter, &update)
	mu.Unlock()
	if err != nil {
		return nil, err
	}

	// Clear data value in cache memory for courses
	wg.Add(2)
	go func() {
		defer wg.Done()
		lessonsCache.Clear()
	}()

	// clear data value with id courseID in cache
	go func() {
		defer wg.Done()
		lessonCache.Remove(lesson.ID.Hex())
	}()
	wg.Wait()

	return data, err
}

// DeleteOne xóa bài học (lesson) theo ID
// Hàm này nhận đầu vào là lessonID và trả về kết quả sau khi xóa
// Nếu có lỗi xảy ra trong quá trình xử lý, hệ thống sẽ trả về lỗi và dừng chương trình
func (l *lessonRepository) DeleteOne(ctx context.Context, lessonID string) error {
	collectionLesson := l.database.Collection(l.collectionLesson)
	collectionUnit := l.database.Collection(l.collectionUnit)

	objID, err := primitive.ObjectIDFromHex(lessonID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	filterInUnit := bson.M{"lesson_id": objID}

	exist, err := collectionUnit.CountDocuments(ctx, filterInUnit)
	if err != nil {
		return err
	}
	if exist > 0 {
		return errors.New(`lesson can not remove`)
	}

	count, err := collectionLesson.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`lesson is removed`)
	}

	_, err = collectionLesson.DeleteOne(ctx, filter)
	// clear data value with courseID in cache
	wg.Add(2)
	go func() {
		defer wg.Done()
		lessonCache.Remove(lessonID)
	}()

	// clear data value with detail in cache due to decrease num
	go func() {
		defer wg.Done()
		detailCache.Clear()
	}()

	wg.Wait()
	return err
}

// Statistics truy vấn thống kê về các bài học (số lượng từ vựng, đơn vị).
// Hàm này không nhận tham số đầu vào và trả về thông tin thống kê.
// Nếu có lỗi xảy ra trong quá trình thống kê, hàm sẽ trả về lỗi và dừng chương trình.
func (l *lessonRepository) Statistics(ctx context.Context) (lesson_domain.Statistics, error) {
	// Khởi tạo một channel để lưu kết quả thống kê
	statisticsCh := make(chan lesson_domain.Statistics, 1)
	// Sử dụng waitGroup để chờ tất cả các goroutine hoàn thành
	wg.Add(1)
	// Bắt đầu một Goroutine để truy vấn dữ liệu bài học từ cache (nếu có)
	go func() {
		defer wg.Done()
		data, found := statisticsCache.Get("statistics")
		if found {
			statisticsCh <- data
			return
		}
	}()

	// Goroutine để đóng channel khi tất cả công việc hoàn thành
	go func() {
		defer close(statisticsCh)
		wg.Wait()
	}()

	// Nhận giá trị từ channel statisticsCh
	statisticsData := <-statisticsCh
	// Kiểm tra nếu statisticsData không null
	// Nếu không, trả về giá trị đó
	if !internal.IsZeroValue(statisticsData) {
		return statisticsData, nil
	}

	// Khởi tạo các bộ sưu tập
	collectionUnit := l.database.Collection(l.collectionUnit)
	collectionVocabulary := l.database.Collection(l.collectionVocabulary)
	collectionLesson := l.database.Collection(l.collectionLesson)

	// Đếm số lượng đơn vị
	countUnit, err := collectionUnit.CountDocuments(ctx, bson.D{})
	if err != nil {
		return lesson_domain.Statistics{}, err
	}

	// Đếm số lượng từ vựng
	countVocabulary, err := collectionVocabulary.CountDocuments(ctx, bson.D{})
	if err != nil {
		return lesson_domain.Statistics{}, err
	}

	// Đếm tổng số lượng bài học
	count, err := collectionLesson.CountDocuments(ctx, bson.D{})
	if err != nil {
		return lesson_domain.Statistics{}, err
	}

	// Tạo cấu trúc Thống kê với dữ liệu đếm
	statistics := lesson_domain.Statistics{
		Total:           count,
		CountUnit:       countUnit,
		CountVocabulary: countVocabulary,
	}

	// Đặt cache memory với dữ liệu cần thiết trong 5 phút
	statisticsCache.Set("statistics", statistics, 5*time.Minute)
	return statistics, nil
}

// countLessonsByCourseID counts the number of lessons associated with a course.
func (l *lessonRepository) countUnitsByLessonsID(ctx context.Context, lessonID primitive.ObjectID) (int32, error) {
	collectionUnit := l.database.Collection(l.collectionUnit)

	filter := bson.M{"lesson_id": lessonID}
	count, err := collectionUnit.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}

func (l *lessonRepository) countVocabularyByLessonID(ctx context.Context, lessonID primitive.ObjectID) (int32, error) {
	collectionVocabulary := l.database.Collection(l.collectionVocabulary)

	// Phần pipeline aggregation để nối các collection và đếm số lượng từ vựng
	pipeline := []bson.M{
		// Nối collection Vocabulary với collection Unit
		{"$lookup": bson.M{
			"from":         "unit",
			"localField":   "unit_id",
			"foreignField": "_id",
			"as":           "unit",
		}},
		{"$unwind": "$unit"},

		// Lọc các từ vựng thuộc về khóa học được cung cấp
		{"$match": bson.M{"unit.lesson_id": lessonID}},
		// Đếm số lượng từ vựng
		{"$count": "totalVocabulary"},
	}

	// Thực hiện aggregation
	cursor, err := collectionVocabulary.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var result struct {
		TotalVocabulary int32 `bson:"totalVocabulary"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	return result.TotalVocabulary, nil
}

// CountUnitsByLessonID đếm số lượng unit trong lesson dựa trên lessonID
func (l *lessonRepository) countUnitsByLessonID(ctx context.Context, lessonID primitive.ObjectID) (int64, error) {
	collectionLesson := l.database.Collection(l.collectionLesson)

	filter := bson.M{"lesson_id": lessonID}
	count, err := collectionLesson.CountDocuments(ctx, filter)
	return count, err
}

// getLastLesson lấy unit cuối cùng từ collection
func (l *lessonRepository) getLastLesson(ctx context.Context) (*lesson_domain.Lesson, error) {
	collectionLesson := l.database.Collection(l.collectionLesson)
	findOptions := options.FindOne().SetSort(bson.D{{"_id", -1}})

	var lesson lesson_domain.Lesson
	err := collectionLesson.FindOne(ctx, bson.D{}, findOptions).Decode(&lesson)
	if err != nil {
		return nil, err
	}

	return &lesson, nil
}
