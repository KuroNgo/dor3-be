package course_repository

import (
	course_domain "clean-architecture/domain/course"
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

type courseRepository struct {
	database             *mongo.Database
	collectionCourse     string
	collectionLesson     string
	collectionUnit       string
	collectionVocabulary string
}

// NewCourseRepository hàm khởi tạo (constructor) để khởi tạo instance của struct
func NewCourseRepository(db *mongo.Database, collectionCourse string, collectionLesson string, collectionUnit string, collectionVocabulary string) course_domain.ICourseRepository {
	return &courseRepository{
		database:             db,
		collectionCourse:     collectionCourse,
		collectionLesson:     collectionLesson,
		collectionUnit:       collectionUnit,
		collectionVocabulary: collectionVocabulary,
	}
}

var (
	courseCache     = cache.NewTTL[string, course_domain.CourseResponse]()
	coursesCache    = cache.NewTTL[string, []course_domain.CourseResponse]()
	detailCache     = cache.NewTTL[string, course_domain.DetailForManyResponse]()
	statisticsCache = cache.NewTTL[string, course_domain.Statistics]()

	wg sync.WaitGroup
	mu sync.Mutex

	// Khởi tạo channel để luu trữ lỗi
	errCh   = make(chan error)
	courses []course_domain.CourseResponse
)

// FetchByID lấy khóa học (course) theo ID
// Hàm này nhận đầu vào (input) là courseID và trả về một bài học làm khóa và nội dung cuủa bài học tương ứng làm giá trị
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả về với các kết quả đã lấy được
func (c *courseRepository) FetchByID(ctx context.Context, courseID string) (course_domain.CourseResponse, error) {
	// Khởi tạo channel để lưu trữ kết quả lesson
	courseCh := make(chan course_domain.CourseResponse)

	wg.Add(1)
	// Khởi động một goroutine cho tìm dữ liệu detail trong cache memory
	go func() {
		defer wg.Done()
		value, found := courseCache.Get(courseID)
		if found {
			courseCh <- value
			return
		}
	}()

	// Goroutine để đóng các channel khi tất cả các công việc hoàn thành
	go func() {
		defer close(courseCh)
		wg.Wait()
	}()

	// Gán giá trị từ channel
	courseData := <-courseCh
	// kiểm tra dữ liệu Data có rỗng hay không,
	// nếu không sẽ trả về dữ lệu trong cache vừa tìm được
	// Ngược lại, sẽ thực hiện quy trình tìm
	if !internal.IsZeroValue(courseData) {
		return courseData, nil
	}

	// khởi tạo đối tượng collection, ở đây là course
	collectionCourse := c.database.Collection(c.collectionCourse)
	// Thực hiện chuyển đổi courseID từ string sang primitive.ObjectID
	idCourse, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return course_domain.CourseResponse{}, err
	}
	// Lấy dữ liệu courseID vừa chuyển đổi, thực hiện tìm kiếm theo id
	filter := bson.M{"_id": idCourse}

	var course course_domain.CourseResponse
	// Thực hiện tìm kiếm course theo id
	err = collectionCourse.FindOne(ctx, filter).Decode(&course)
	if err != nil {
		return course_domain.CourseResponse{}, err
	}

	countLessonCh := make(chan int32)
	countVocabularyCh := make(chan int32)

	// Goroutine để thực hiên đếm số lượng lesson trong lesson của course
	go func() {
		defer close(countLessonCh)
		countLesson, err := c.countLessonsByCourseID(ctx, course.Id)
		if err != nil {
			errCh <- err
			return
		}
		countLessonCh <- countLesson
	}()

	// Goroutine để thực hiên đếm số lượng vocabulary trong lesson
	go func() {
		defer close(countVocabularyCh)
		countVocabulary, err := c.countVocabularyByCourseID(ctx, course.Id)
		if err != nil {
			errCh <- err
			return
		}
		countVocabularyCh <- countVocabulary
	}()

	// Channel gửi giá trị, sau đó lesson sẽ nhận giá trị tương ứng
	countLesson := <-countLessonCh
	countVocab := <-countVocabularyCh

	course.CountVocabulary = countVocab
	course.CountLesson = countLesson

	// Thiết lập Set cache memory với dữ liệu cần thiết với thơi gian là 5 phút
	courseCache.Set(courseID, course, 5*time.Minute)

	select {
	// Nếu có lỗi, sẽ thực hiện trả về lỗi
	case err = <-errCh:
		return course_domain.CourseResponse{}, err
	// Ngược lại, sẽ trả về giá trị
	default:
		return course, nil
	}
}

// FetchManyForEachCourse lấy tất cả khóa học (course) cùng một lúc (concurrency).
// Hàm này nhận vào số trang (page) và trả về một mảng khóa học làm khóa và nội dung của bài học tương ứng làm giá trị.
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả với các kết quả đã lấy được
// FIXME: thực hiện gắn lỗi vào channel giúp tối ưu hóa xử lý
func (c *courseRepository) FetchManyForEachCourse(ctx context.Context, page string) ([]course_domain.CourseResponse, course_domain.DetailForManyResponse, error) {
	// buffer channel with increase performance
	// note: use buffer have target
	// Khởi tạo channel để lưu trữ kết quả course
	coursesCh := make(chan []course_domain.CourseResponse, 5)
	// Khởi tạo channel để lưu trữ kết quả detail
	detailCh := make(chan course_domain.DetailForManyResponse, 1)
	// Sử dụng WaitGroup để đợi tất cả các goroutine hoàn thành
	wg.Add(2)
	// Khởi động một goroutine cho tìm dữ liệu lesson trong cache memory
	go func() {
		defer wg.Done()
		data, found := coursesCache.Get(page)
		if found {
			coursesCh <- data
			return
		}
	}()

	// Khởi động một goroutine cho tìm dữ liệu detail trong cache memory
	go func() {
		defer wg.Done()
		detailData, foundDetail := detailCache.Get("detail")
		if foundDetail {
			detailCh <- detailData
			return
		}
	}()

	// Goroutine để đóng các channel khi tất cả các công việc hoàn thành
	go func() {
		defer close(coursesCh)
		defer close(detailCh)
		wg.Wait()
	}()

	// Gán giá trị từ channel
	courseData := <-coursesCh
	responseData := <-detailCh

	// kiểm tra dữ liệu Data có rỗng hay không,
	// nếu không sẽ trả về dữ lệu trong cache vừa tìm được
	// Ngược lại, sẽ thực hiện quy trình tìm
	if !internal.IsZeroValue(courseData) && !internal.IsZeroValue(responseData) {
		return courseData, responseData, nil
	}

	// khởi tạo đối tượng collection, ở đây là course
	collectionCourse := c.database.Collection(c.collectionCourse)

	// thực hiện phân trang
	// lấy số trang từ client
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, errors.New("invalid page number")
	}
	// tối đa dữ liệu gửi đến ở mỗi trang được yêu cầu
	perPage := 5
	// nếu dữ liệu nhỏ hơn 1 sẽ skip
	skip := (pageNumber - 1) * perPage
	// thực hiện các yêu cầu đã neu ở trên
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	// đếm số lượng lesson có trong project
	count, err := collectionCourse.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, err
	}

	// đếm số lượng trang
	totalPages := (count + int64(perPage) - 1) / int64(perPage)

	// thực hiện tìm kiếm theo điêu kiện options
	cursor, err := collectionCourse.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, err
	}
	// thực hiện close cursor nếu hàm đã hoàn thành hoặc bị lỗi
	defer func() {
		err = cursor.Close(ctx)
		if err != nil {
			return
		}
	}()

	wg.Add(1)
	// Khởi động một goroutine cho mỗi cursor
	// TODO: Xử lý tìm từng dữ liệu liên quan đến course bao gồm các thông tin cơ bản và các thông tin khác
	go func() {
		defer wg.Done()
		for cursor.Next(ctx) {
			// chuyển đổi sang JSON cho course
			var course course_domain.CourseResponse
			if err = cursor.Decode(&course); err != nil {
				errCh <- err
				return
			}

			wg.Add(1)
			// Goroutine giúp lấy kết quả đếm lesson, vocabulary cho mỗi course
			go func(course2 course_domain.CourseResponse) {
				defer wg.Done()
				countLesson, err := c.countLessonsByCourseID(ctx, course2.Id)
				if err != nil {
					errCh <- err
					return
				}

				countVocab, err := c.countVocabularyByCourseID(ctx, course2.Id)
				if err != nil {
					errCh <- err
					return
				}
				course2.CountVocabulary = countVocab
				course2.CountLesson = countLesson
				courses = append(courses, course2)
			}(course)
		}
	}()
	wg.Wait()

	// Channel để thu thập kết quả thống kê
	statisticsCh := make(chan course_domain.Statistics)
	go func() {
		defer close(statisticsCh)
		statistics, _ := c.Statistics(ctx)
		statisticsCh <- statistics
	}()
	statistics := <-statisticsCh

	// Luu tất cả thông tin liên quan vào detail
	detail := course_domain.DetailForManyResponse{
		CountCourse: count,
		Page:        totalPages,
		Statistics:  statistics,
		CurrentPage: pageNumber,
	}

	// Thiết lập Set cache memory với dữ liệu cần thiết với thơi gian là 5 phút
	coursesCache.Set(page, courses, 5*time.Minute)
	detailCache.Set("detail", detail, 5*time.Minute)

	// Thu thập kết quả
	select {
	// Nếu có lỗi, sẽ thực hiện trả về lỗi
	case err = <-errCh:
		return nil, course_domain.DetailForManyResponse{}, err
	// Ngược lại, sẽ trả về giá trị
	default:
		return courses, detail, nil
	}
}

// UpdateOne cập nhật khóa học (course) theo đối tượng course
// Hàm nhận tham số là đối tượng Course. Nếu thành công sẽ trả về thông tin cập nhật (không phải thông tin đối tượng).
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả về với kết quả đã lấy được
// Hàm có sử dụng concurrency, giúp xử lý các tác vụ về người dùng hiệu quả
func (c *courseRepository) UpdateOne(ctx context.Context, course *course_domain.Course) (*mongo.UpdateResult, error) {
	// khởi tạo đối tượng collection, ở đây là course
	collectionCourse := c.database.Collection(c.collectionCourse)

	// Thực hiện tìm kiếm theo id
	filter := bson.D{{Key: "_id", Value: course.Id}}
	// Thực hiện cập nhật đối tương theo các trường cho trước
	update := bson.M{
		"$set": bson.M{
			"name":        course.Name,
			"description": course.Description,
			"updated_at":  course.UpdatedAt,
			"who_updated": course.WhoUpdated,
		},
	}

	// Khóa lock giúp bảo vệ course
	mu.Lock()
	data, err := collectionCourse.UpdateOne(ctx, filter, &update)
	// Mở lock khi thực thi xong hành động update
	mu.Unlock()
	if err != nil {
		return nil, err
	}

	// Clear data value in cache memory for courses
	wg.Add(3)
	go func() {
		defer wg.Done()
		coursesCache.Clear()
	}()

	// clear data value with id courseID in cache
	go func() {
		defer wg.Done()
		courseCache.Remove(course.Id.Hex())
	}()
	wg.Wait()

	go func() {
		defer wg.Done()
		statisticsCache.Clear()
	}()
	return data, nil
}

// CreateOne khởi tạo khóa học (course) theo đối tượng Course
// Hàm này nhận tham số là một đối tượng và trả về kết quả thông tin xử lý (không phải là thông tin của đối tượng đó)
// Nếu có lỗi xảy ra trong quá trình xử lý, lỗi sẽ được trả về với kết quả đã lấy được và dừng chương trình
func (c *courseRepository) CreateOne(ctx context.Context, course *course_domain.Course) error {
	// khởi tạo đối tượng collection, ở đây là course
	collectionCourse := c.database.Collection(c.collectionCourse)

	// Thực hiện tìm kiếm theo name để kiểm tra có dữ liệu trùng không
	filter := bson.M{"name": course.Name}
	count, err := collectionCourse.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the course name already exists")
	}

	_, err = collectionCourse.InsertOne(ctx, course)

	if err != nil {
		return err
	}

	// Clear data value in cache memory
	wg.Add(3)
	go func() {
		defer wg.Done()
		coursesCache.Clear()
	}()

	// clear data value in cache memory due to increase num
	go func() {
		defer wg.Done()
		detailCache.Clear()
	}()

	// clear data value in cache memory due to increase num
	go func() {
		defer wg.Done()
		statisticsCache.Clear()
	}()
	wg.Wait()

	return nil
}

// DeleteOne xóa khóa học (course) theo ID
// Hàm này nhận đầu vào là courseID và trả về kết quả sau khi xóa
// Nếu có lỗi xảy ra trong quá trình xử lý, hệ thống sẽ trả về lỗi và dừng chương trình
func (c *courseRepository) DeleteOne(ctx context.Context, courseID string) error {
	collectionCourse := c.database.Collection(c.collectionCourse)

	// Convert courseID string to ObjectID
	objID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return err
	}

	// Check if any lesson is associated with the course
	countFK, err := c.countLessonsByCourseID(ctx, objID)
	if err != nil {
		return err
	}
	if countFK > 0 {
		return errors.New("the course cannot be deleted because it is associated with lessons")
	}

	filter := bson.M{"_id": objID}
	// Khóa lock giúp bảo vệ course
	mu.Lock()
	_, err = collectionCourse.DeleteOne(ctx, filter)
	// Mở lock khi thực thi xong hành động delete
	mu.Unlock()
	if err != nil {
		return err
	}

	// clear data value with courseID in cache
	wg.Add(3)
	go func() {
		defer wg.Done()
		courseCache.Remove(courseID)
	}()

	// clear data value with detail in cache due to decrease num
	go func() {
		defer wg.Done()
		detailCache.Clear()
	}()

	// clear data value with detail in cache due to decrease num
	go func() {
		defer wg.Done()
		statisticsCache.Clear()
	}()
	wg.Wait()

	return nil
}

// Statistics thống kê khóa học (có bao nhiêu bài học, số lượng từ vựng, unit)
// Hàm này không nhận đầu vào (input), trả về thông tin thống kê
// Nếu có lỗi xảy ra trong quá trình thống kê, lỗi đó sẽ được trả về và dừng chương trình
func (c *courseRepository) Statistics(ctx context.Context) (course_domain.Statistics, error) {
	// Khởi tạo channel để lưu kết quả statistics
	statisticsCh := make(chan course_domain.Statistics, 1)
	// Sử dụng waitGroup để đợi tất cả goroutine hoàn thành
	wg.Add(1)
	// Khởi động Goroutine giúp tìm dữ liệu lesson
	// theo id trong cache (đã từng tìm lessonID này hay chưa)
	go func() {
		defer wg.Done()
		data, found := statisticsCache.Get("statistics")
		if found {
			statisticsCh <- data
			return
		}
	}()

	// Goroutine để đóng các channel khi tất cả các công việc hoàn thành
	go func() {
		defer close(statisticsCh)
		wg.Wait()
	}()

	// Channel gửi giá trị cho biến statisticsData
	statisticsData := <-statisticsCh
	// Kiểm tra giá trị statisticsData có null ?
	// Nếu không thì sẽ thực hiện trả về giá trị
	if !internal.IsZeroValue(statisticsData) {
		return statisticsData, nil
	}

	collectionCourse := c.database.Collection(c.collectionCourse)
	count, err := collectionCourse.CountDocuments(ctx, bson.D{})
	if err != nil {
		return course_domain.Statistics{}, err
	}

	statistics := course_domain.Statistics{
		Total: count,
	}

	// Thiết lập Set cache memory với dữ liệu cần thiết với thơi gian là 5 phút
	statisticsCache.Set("statistics", statistics, 5*time.Minute)
	return statistics, nil
}

// countLessonsByCourseID counts the number of lessons associated with a course.
func (c *courseRepository) countLessonsByCourseID(ctx context.Context, courseID primitive.ObjectID) (int32, error) {
	collectionLesson := c.database.Collection(c.collectionLesson)

	filter := bson.M{"course_id": courseID}
	count, err := collectionLesson.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}

// countVocabularyByCourseID counts the number of vocabularies associated with a course.
func (c *courseRepository) countVocabularyByCourseID(ctx context.Context, courseID primitive.ObjectID) (int32, error) {
	collectionVocabulary := c.database.Collection(c.collectionVocabulary)

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
		// Nối collection Unit với collection Lesson
		{"$lookup": bson.M{
			"from":         "lesson",
			"localField":   "unit.lesson_id",
			"foreignField": "_id",
			"as":           "lesson",
		}},
		{"$unwind": "$lesson"},
		// Lọc các từ vựng thuộc về khóa học được cung cấp
		{"$match": bson.M{"lesson.course_id": courseID}},
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
