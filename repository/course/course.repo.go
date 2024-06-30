package course_repository

import (
	course_domain "clean-architecture/domain/course"
	"clean-architecture/internal"
	"clean-architecture/internal/cache/memory"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
	"time"
)

type courseRepository struct {
	database                *mongo.Database
	collectionCourse        string
	collectionCourseProcess string
	collectionLesson        string
	collectionUnit          string
	collectionVocabulary    string
}

// NewCourseRepository hàm khởi tạo (constructor) để khởi tạo instance của struct
func NewCourseRepository(db *mongo.Database, collectionCourse string, collectionCourseProcess string, collectionLesson string, collectionUnit string, collectionVocabulary string) course_domain.ICourseRepository {
	return &courseRepository{
		database:                db,
		collectionCourse:        collectionCourse,
		collectionCourseProcess: collectionCourseProcess,
		collectionLesson:        collectionLesson,
		collectionUnit:          collectionUnit,
		collectionVocabulary:    collectionVocabulary,
	}
}

var (
	courseCache             = memory.NewTTL[string, course_domain.CourseResponse]()
	coursesCache            = memory.NewTTL[string, []course_domain.CourseResponse]()
	coursePrimOIDCache      = memory.NewTTL[string, primitive.ObjectID]()
	coursesUserProcessCache = memory.NewTTL[string, []course_domain.CourseProcess]()
	courseUserProcessCache  = memory.NewTTL[string, course_domain.CourseProcess]()
	detailCourseCache       = memory.NewTTL[string, course_domain.DetailForManyResponse]()
	statisticsCache         = memory.NewTTL[string, course_domain.Statistics]()

	wg           sync.WaitGroup
	mu           sync.Mutex
	isProcessing bool
)

// FetchByIDInUser retrieves the course process data for a specific user and course.
// It first checks the cache for the data. If not found, it queries the database and,
// if necessary, initializes the course process data for the user.
//
// Parameters:
// - ctx: The context for managing request deadlines, cancellation signals, and other request-scoped values.
// - userID: The ObjectID of the user whose course process data is being retrieved.
// - courseID: The string representation of the ObjectID of the course.
//
// Returns:
// - course_domain.CourseProcess: The course process data for the user.
// - error: Any error encountered during the operation.
func (c *courseRepository) FetchByIDInUser(ctx context.Context, userID primitive.ObjectID, courseID string) (course_domain.CourseProcess, error) {
	// Create a channel to capture errors and defer its closure
	errCh := make(chan error)
	defer close(errCh)

	// Create a channel to capture the course process data
	courseUserProcessCh := make(chan course_domain.CourseProcess)

	// Add a goroutine to the wait group to fetch data from the cache
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Try to get the data from cache
		data, found := courseUserProcessCache.Get(userID.Hex() + courseID)
		if found {
			// Send the data to the channel if found
			courseUserProcessCh <- data
		}
	}()

	// Close the courseUserProcessCh channel after all goroutines are done
	go func() {
		defer close(courseUserProcessCh)
		wg.Wait()
	}()

	// Wait for data from the courseUserProcessCh channel
	courseUserProcessData := <-courseUserProcessCh
	if !internal.IsZeroValue(courseUserProcessData) {
		// Return the data if it is not zero value
		return courseUserProcessData, nil
	}

	// Get the course and course process collections from the database
	collectionCourse := c.database.Collection(c.collectionCourse)
	collectionCourseProcess := c.database.Collection(c.collectionCourseProcess)

	// Convert courseID from string to ObjectID
	idCourse, _ := primitive.ObjectIDFromHex(courseID)
	filterCourseProcessByUser := bson.M{"user_id": userID, "course_id": idCourse}

	// Count the total number of courses
	countCourse, err := collectionCourse.CountDocuments(ctx, bson.D{})
	if err != nil {
		return course_domain.CourseProcess{}, err
	}

	// Count the number of CourseProcess documents for the user
	count, err := collectionCourseProcess.CountDocuments(ctx, filterCourseProcessByUser)
	if err != nil {
		return course_domain.CourseProcess{}, err
	}

	// If the user does not have CourseProcess documents for all courses, initialize them
	var courseUserProcess course_domain.CourseProcess
	if count < countCourse {
		// Find all courses
		cursorCourse, err := collectionCourse.Find(ctx, bson.D{})
		if err != nil {
			return course_domain.CourseProcess{}, err
		}
		defer func(cursorCourse *mongo.Cursor, ctx context.Context) {
			err := cursorCourse.Close(ctx)
			if err != nil {
				// Send any errors to the error channel
				errCh <- err
				return
			}
		}(cursorCourse, ctx)

		// Iterate over all courses
		for cursorCourse.Next(ctx) {
			var course course_domain.Course
			if err = cursorCourse.Decode(&course); err != nil {
				return course_domain.CourseProcess{}, err
			}

			// Add a goroutine to the wait group to process each course
			wg.Add(1)
			go func(course course_domain.Course) {
				defer wg.Done()
				// Initialize the course process for the user
				courseProcess := course_domain.CourseProcess{
					CourseID:   course.Id,
					UserID:     userID,
					IsComplete: 0,
				}

				// Check if a CourseProcess document already exists
				filter := bson.M{"course_id": course.Id, "user_id": userID}
				countCourseChild, err := collectionCourseProcess.CountDocuments(ctx, filter)
				if err != nil {
					// Send any errors to the error channel
					errCh <- err
					return
				}

				// If no document exists, insert the new course process
				if countCourseChild == 0 {
					_, err = collectionCourseProcess.InsertOne(ctx, &courseProcess)
					if err != nil {
						log.Println("Error inserting course process:", err)
						// Send any errors to the error channel
						errCh <- err
						return
					}
				}
			}(course)
			wg.Wait()
		}

		// Find the CourseProcess document for the user with pagination
		err = collectionCourseProcess.FindOne(ctx, filterCourseProcessByUser).Decode(&courseUserProcess)
		if err != nil {
			return course_domain.CourseProcess{}, err
		}
	}

	// Find the CourseProcess document for the user
	err = collectionCourseProcess.FindOne(ctx, filterCourseProcessByUser).Decode(&courseUserProcess)
	if err != nil {
		return course_domain.CourseProcess{}, err
	}

	// Cache the CourseProcess data for 5 minutes
	courseUserProcessCache.Set(userID.Hex()+courseID, courseUserProcess, 5*time.Minute)

	// Check if there were any errors in the error channel
	select {
	case err = <-errCh:
		return course_domain.CourseProcess{}, err
	default:
		return courseUserProcess, nil
	}
}

// FetchManyInUser retrieves multiple course processes for a specific user with pagination support.
// It first checks the cache for the data. If not found, it queries the database and,
// if necessary, initializes the course process data for the user.
//
// Parameters:
// - ctx: The context for managing request deadlines, cancellation signals, and other request-scoped values.
// - userID: The ObjectID of the user whose course process data is being retrieved.
// - page: The page number for pagination.
//
// Returns:
// - []course_domain.CourseProcess: A slice of course process data for the user.
// - course_domain.DetailForManyResponse: Detailed response including statistics and pagination info.
// - error: Any error encountered during the operation.
func (c *courseRepository) FetchManyInUser(ctx context.Context, userID primitive.ObjectID, page string) ([]course_domain.CourseProcess, course_domain.DetailForManyResponse, error) {
	errCh := make(chan error)
	defer close(errCh)

	courseUserProcessCh := make(chan []course_domain.CourseProcess)
	detailCh := make(chan course_domain.DetailForManyResponse)

	wg.Add(2)
	go func() {
		defer wg.Done()
		data, found := coursesUserProcessCache.Get(userID.Hex())
		if found {
			courseUserProcessCh <- data
		}
	}()

	go func() {
		defer wg.Done()
		data, found := detailCourseCache.Get(userID.Hex() + "detail")
		if found {
			detailCh <- data
		}
	}()

	go func() {
		defer close(detailCh)
		defer close(courseUserProcessCh)
		wg.Wait()
	}()

	courseUserProcessData := <-courseUserProcessCh
	detailData := <-detailCh
	if !internal.IsZeroValue(courseUserProcessData) && !internal.IsZeroValue(detailData) {
		return courseUserProcessData, detailData, nil
	}

	collectionCourse := c.database.Collection(c.collectionCourse)
	collectionCourseProcess := c.database.Collection(c.collectionCourseProcess)

	filterCourseProcessByUser := bson.M{"user_id": userID}

	// Thực hiện phân trang
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, errors.New("invalid page number")
	}
	perPage := 5
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	// Đếm số lượng khóa học trong collection 'courses'
	countCourse, err := collectionCourse.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, err
	}

	// Tính toán tổng số trang dựa trên số lượng khóa học và số khóa học mỗi trang
	totalPages := (countCourse + int64(perPage) - 1) / int64(perPage)

	// Đếm số lượng CourseProcess của người dùng
	count, err := collectionCourseProcess.CountDocuments(ctx, filterCourseProcessByUser)
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, err
	}

	var coursesProcess []course_domain.CourseProcess
	// Nếu không có CourseProcess cho người dùng, khởi tạo chúng
	if count < countCourse {
		cursorCourse, err := collectionCourse.Find(ctx, bson.D{})
		if err != nil {
			return nil, course_domain.DetailForManyResponse{}, err
		}
		defer func(cursorCourse *mongo.Cursor, ctx context.Context) {
			err := cursorCourse.Close(ctx)
			if err != nil {
				errCh <- err
				return
			}
		}(cursorCourse, ctx)

		for cursorCourse.Next(ctx) {
			var course course_domain.Course
			if err = cursorCourse.Decode(&course); err != nil {
				return nil, course_domain.DetailForManyResponse{}, err
			}

			wg.Add(1)
			go func(course course_domain.Course) {
				defer wg.Done()
				courseProcess := course_domain.CourseProcess{
					CourseID:   course.Id,
					UserID:     userID,
					IsComplete: 0,
				}

				// Thực hiện tìm kiếm theo name để kiểm tra có dữ liệu trùng không
				filter := bson.M{"course_id": course.Id, "user_id": userID}
				countCourseChild, err := collectionCourseProcess.CountDocuments(ctx, filter)
				if err != nil {
					errCh <- err
					return
				}

				if countCourseChild == 0 {
					_, err = collectionCourseProcess.InsertOne(ctx, &courseProcess)
					if err != nil {
						log.Println("Error inserting course process:", err)
						errCh <- err
						return
					}
				}
			}(course)
		}
		wg.Wait()

		// Tìm các CourseProcess của người dùng với phân trang
		cursor, err := collectionCourseProcess.Find(ctx, filterCourseProcessByUser, findOptions)
		if err != nil {
			return nil, course_domain.DetailForManyResponse{}, err
		}
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			err := cursor.Close(ctx)
			if err != nil {
				errCh <- err
				return
			}
		}(cursor, ctx)
	}

	// Tìm các CourseProcess của người dùng với phân trang
	cursor, err := collectionCourseProcess.Find(ctx, filterCourseProcessByUser, findOptions)
	if err != nil {
		return nil, course_domain.DetailForManyResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	// Đọc dữ liệu từ cursor và thêm vào slice coursesProcess
	for cursor.Next(ctx) {
		var courseProcess course_domain.CourseProcess
		if err := cursor.Decode(&courseProcess); err != nil {
			return nil, course_domain.DetailForManyResponse{}, err
		}
		mu.Lock()
		coursesProcess = append(coursesProcess, courseProcess)
		mu.Unlock()
	}

	if err := cursor.Err(); err != nil {
		return nil, course_domain.DetailForManyResponse{}, err
	}

	// Lấy thống kê cho detail response
	statistics, _ := c.Statistics(ctx, filterCourseProcessByUser)
	detail := course_domain.DetailForManyResponse{
		Statistics:  statistics,
		Page:        totalPages,
		CurrentPage: pageNumber,
		CountCourse: countCourse,
	}

	coursesUserProcessCache.Set(userID.Hex(), coursesProcess, 5*time.Minute)
	detailCourseCache.Set(userID.Hex()+"detail", detail, 5*time.Minute)

	select {
	case err = <-errCh:
		return nil, course_domain.DetailForManyResponse{}, err
	default:
		return coursesProcess, detail, nil
	}
}

func (c *courseRepository) UpdateCompleteInUser(ctx context.Context, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	//TODO implement me
	panic("implement me")
}

// FetchByIDInAdmin lấy khóa học (course) theo ID
// Hàm này nhận đầu vào (input) là courseID và trả về một bài học làm khóa và nội dung cuủa bài học tương ứng làm giá trị
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả về với các kết quả đã lấy được
func (c *courseRepository) FetchByIDInAdmin(ctx context.Context, courseID string) (course_domain.CourseResponse, error) {
	// Khởi tạo channel để luu trữ lỗi
	errCh := make(chan error, 1)
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

// FetchManyForEachCourseInAdmin lấy tất cả khóa học (course) cùng một lúc (concurrency).
// Hàm này nhận vào số trang (page) và trả về một mảng khóa học làm khóa và nội dung của bài học tương ứng làm giá trị.
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả với các kết quả đã lấy được
func (c *courseRepository) FetchManyForEachCourseInAdmin(ctx context.Context, page string) ([]course_domain.CourseResponse, course_domain.DetailForManyResponse, error) {
	// Khởi tạo channel để luu trữ lỗi
	errCh := make(chan error, 1)
	// Khởi tạo channel để lưu trữ kết quả course
	coursesCh := make(chan []course_domain.CourseResponse, 1)
	// Khởi tạo channel để lưu trữ kết quả detail
	detailCh := make(chan course_domain.DetailForManyResponse, 1)
	// Sử dụng WaitGroup để đợi tất cả các goroutine hoàn thành
	wg.Add(2)
	// Khởi động một goroutine cho tìm dữ liệu lesson trong cache memory
	go func() {
		defer wg.Done()
		defer close(coursesCh)
		data, found := coursesCache.Get(page)
		if found {
			coursesCh <- data
			return
		}
	}()

	// Khởi động một goroutine cho tìm dữ liệu detail trong cache memory
	go func() {
		defer wg.Done()
		defer close(detailCh)
		detailData, foundDetail := detailCourseCache.Get("detail")
		if foundDetail {
			detailCh <- detailData
			return
		}
	}()

	wg.Wait()

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
	filter := bson.M{}
	cursor, err := collectionCourse.Find(ctx, filter, findOptions)
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

	var courses []course_domain.CourseResponse
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
		statistics, _ := c.Statistics(ctx, filter)
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
	detailCourseCache.Set("detail", detail, 5*time.Minute)

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

// FindCourseIDByCourseNameInAdmin retrieves the course ID (courseid) based on the given course name.
// This function accepts courseName as a parameter and returns a primitive.ObjectID.
// If an error occurs during data retrieval, the error is returned along with any partially retrieved results.
// TODO: This function is intended to be a helper and not used as an API (controller).
func (c *courseRepository) FindCourseIDByCourseNameInAdmin(ctx context.Context, courseName string) (primitive.ObjectID, error) {
	// Channel to receive the course ID
	courseIDCh := make(chan primitive.ObjectID)

	// Add a goroutine to the wait group
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Try to get the course ID from the cache
		data, found := coursePrimOIDCache.Get(courseName)
		if found {
			// Send the cached course ID to the channel
			courseIDCh <- data
			return
		}
	}()

	// Goroutine to wait for other goroutines and then close the channel
	go func() {
		defer close(courseIDCh)
		wg.Wait()
	}()

	// Receive the course ID data from the channel
	courseIDData := <-courseIDCh
	// Check if the course ID data is non-zero
	if !internal.IsZeroValue(courseIDData) {
		return courseIDData, nil
	}

	// Connect to the courses collection in the database
	collectionCourse := c.database.Collection(c.collectionCourse)

	// Filter to find the course by name
	filter := bson.M{"name": courseName}
	// Struct to hold the result
	var data struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	// Find one course matching the filter and decode the result
	err := collectionCourse.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		// Return an error if the course is not found or another error occurs
		return primitive.NilObjectID, err
	}

	// Cache the retrieved course ID
	coursePrimOIDCache.Set(courseName, data.Id, 5*time.Minute)
	// Return the retrieved course ID
	return data.Id, nil
}

// UpdateOneInAdmin cập nhật khóa học (course) theo đối tượng course
// Hàm nhận tham số là đối tượng Course. Nếu thành công sẽ trả về thông tin cập nhật (không phải thông tin đối tượng).
// Nếu có lỗi xảy ra trong quá trình lấy dữ liệu, lỗi đó sẽ được trả về với kết quả đã lấy được
// Hàm có sử dụng concurrency, giúp xử lý các tác vụ về người dùng hiệu quả
func (c *courseRepository) UpdateOneInAdmin(ctx context.Context, course *course_domain.Course) (*mongo.UpdateResult, error) {
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

	// khởi tạo đối tượng collection, ở đây là course
	collectionCourse := c.database.Collection(c.collectionCourse)

	if course.Id == primitive.NilObjectID {
		return nil, errors.New("the course id not nil")
	}

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

	data, err := collectionCourse.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	// Clear data value in cache memory for courses
	wg.Add(4)
	go func() {
		defer wg.Done()
		coursesCache.Clear()
	}()

	// clear data value with id courseID in cache
	go func() {
		defer wg.Done()
		courseCache.Remove(course.Id.Hex())
	}()

	go func() {
		defer wg.Done()
		statisticsCache.Clear()
	}()

	go func() {
		defer wg.Done()
		coursesUserProcessCache.Clear()
	}()
	wg.Wait()

	return data, nil
}

// CreateOneInAdmin khởi tạo khóa học (course) theo đối tượng Course
// Hàm này nhận tham số là một đối tượng và trả về kết quả thông tin xử lý (không phải là thông tin của đối tượng đó)
// Nếu có lỗi xảy ra trong quá trình xử lý, lỗi sẽ được trả về với kết quả đã lấy được và dừng chương trình
func (c *courseRepository) CreateOneInAdmin(ctx context.Context, course *course_domain.Course) error {
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

	// khởi tạo đối tượng collection, ở đây là course
	collectionCourse := c.database.Collection(c.collectionCourse)

	if course.Id == primitive.NilObjectID {
		return errors.New("the course id not nil")
	}

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
	wg.Add(4)
	go func() {
		defer wg.Done()
		coursesCache.Clear()
	}()

	// clear data value in cache memory due to increase num
	go func() {
		defer wg.Done()
		detailCourseCache.Clear()
	}()

	// clear data value in cache memory due to increase num
	go func() {
		defer wg.Done()
		statisticsCache.Clear()
	}()

	go func() {
		defer wg.Done()
		coursesUserProcessCache.Clear()
	}()
	wg.Wait()

	return nil
}

// DeleteOneInAdmin xóa khóa học (course) theo ID
// Hàm này nhận đầu vào là courseID và trả về kết quả sau khi xóa
// Nếu có lỗi xảy ra trong quá trình xử lý, hệ thống sẽ trả về lỗi và dừng chương trình
func (c *courseRepository) DeleteOneInAdmin(ctx context.Context, courseID string) error {
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

	collectionCourse := c.database.Collection(c.collectionCourse)
	collectionCourseProcess := c.database.Collection(c.collectionCourseProcess)

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
	count, err := collectionCourse.CountDocuments(ctx, filter)
	if count == 0 {
		return errors.New("the course id do not exist")
	}

	_, err = collectionCourse.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	filterCourse := bson.M{"course_id": objID}
	_, err = collectionCourseProcess.DeleteOne(ctx, filterCourse)
	if err != nil {
		return err
	}

	// clear data value with courseID in cache
	wg.Add(6)
	go func() {
		defer wg.Done()
		courseCache.Remove(courseID)
	}()

	// clear data value with detail in cache due to decrease num
	go func() {
		defer wg.Done()
		detailCourseCache.Clear()
	}()

	// clear data value with detail in cache due to decrease num
	go func() {
		defer wg.Done()
		statisticsCache.Clear()
	}()

	// clear data value with detail in cache due to decrease num
	go func() {
		defer wg.Done()
		coursePrimOIDCache.Clear()
	}()

	go func() {
		defer wg.Done()
		coursesUserProcessCache.Clear()
	}()

	go func() {
		defer wg.Done()
		courseUserProcessCache.Clear()
	}()
	wg.Wait()

	return nil
}

// Statistics thống kê khóa học (có bao nhiêu bài học, số lượng từ vựng, unit)
// Hàm này không nhận đầu vào (input), trả về thông tin thống kê
// Nếu có lỗi xảy ra trong quá trình thống kê, lỗi đó sẽ được trả về và dừng chương trình
func (c *courseRepository) Statistics(ctx context.Context, countOptions bson.M) (course_domain.Statistics, error) {
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
	count, err := collectionCourse.CountDocuments(ctx, countOptions)
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
