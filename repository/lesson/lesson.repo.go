package lesson_repository

import (
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type lessonRepository struct {
	database         mongo.Database
	collectionLesson string
	collectionCourse string
}

func NewLessonRepository(db mongo.Database, collectionLesson string, collectionCourse string) lesson_domain.ILessonRepository {
	return &lessonRepository{
		database:         db,
		collectionLesson: collectionLesson,
		collectionCourse: collectionCourse,
	}
}
func (l *lessonRepository) FetchByIdCourse(ctx context.Context, idCourse string) (lesson_domain.Response, error) {
	collectionLesson := l.database.Collection(l.collectionLesson)
	collectionCourse := l.database.Collection(l.collectionCourse)

	idCourse2, err := primitive.ObjectIDFromHex(idCourse)
	filter := bson.M{"course_id": idCourse2}

	cursor, err := collectionLesson.Find(ctx, filter)
	if err != nil {
		return lesson_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var lessons []lesson_domain.Lesson
	// Lặp qua các kết quả và giải mã vào slice units
	for cursor.Next(ctx) {
		var lesson lesson_domain.Lesson
		if err = cursor.Decode(&lesson); err != nil {
			return lesson_domain.Response{}, err
		}

		var course course_domain.Course
		err = collectionCourse.FindOne(ctx, bson.M{"_id": idCourse2}).Decode(&course)
		if err != nil {
			return lesson_domain.Response{}, err
		}

		lesson.CourseID = idCourse2

		lessons = append(lessons, lesson)
	}

	// Tạo và trả về phản hồi với dữ liệu units và số lượng tài liệu trong collection bài học
	response := lesson_domain.Response{
		Lesson: lessons,
	}
	return response, nil
}

func (l *lessonRepository) UpdateComplete(ctx context.Context, lessonID string, lesson lesson_domain.Lesson) error {
	//TODO implement me
	panic("implement me")
}

func (l *lessonRepository) FetchMany(ctx context.Context) (lesson_domain.Response, error) {
	collectionLesson := l.database.Collection(l.collectionLesson)
	collectionCourse := l.database.Collection(l.collectionCourse)

	cursor, err := collectionLesson.Find(ctx, bson.D{})
	if err != nil {
		return lesson_domain.Response{}, err
	}
	defer func(cursor mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var lessons []lesson_domain.Lesson
	for cursor.Next(ctx) {
		var lesson lesson_domain.Lesson
		if err = cursor.Decode(&lesson); err != nil {
			return lesson_domain.Response{}, err
		}

		var course course_domain.Course
		err = collectionCourse.FindOne(ctx, bson.M{"_id": lesson.CourseID}).Decode(&course)
		if err != nil {
			return lesson_domain.Response{}, err
		}

		// Gắn tên của course vào lesson
		lesson.CourseID = course.Id

		// Thêm lesson vào slice lessons
		lessons = append(lessons, lesson)
	}
	err = cursor.All(ctx, &lessons)
	lessonRes := lesson_domain.Response{
		Lesson: lessons,
	}

	return lessonRes, err
}

func (l *lessonRepository) UpdateOne(ctx context.Context, lessonID string, lesson lesson_domain.Lesson) error {
	collection := l.database.Collection(l.collectionLesson)
	doc, err := internal.ToDoc(lesson)
	objID, err := primitive.ObjectIDFromHex(lessonID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

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

	filterReference := bson.M{"_id": lesson.CourseID}
	count, err = collectionCourse.CountDocuments(ctx, filterReference)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("the course ID do not exist")
	}

	_, err = collectionLesson.InsertOne(ctx, lesson)
	return nil
}

func (l *lessonRepository) UpsertOne(ctx context.Context, id string, lesson *lesson_domain.Lesson) (*lesson_domain.Lesson, error) {
	collectionLesson := l.database.Collection(l.collectionLesson)
	collectionCourse := l.database.Collection(l.collectionCourse)

	filterReference := bson.M{"_id": lesson.CourseID}
	count, err := collectionCourse.CountDocuments(ctx, filterReference)
	if err != nil {
		return nil, err
	}

	if count <= 0 {
		return nil, errors.New("the course ID do not exist")
	}

	doc, err := internal.ToDoc(lesson)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collectionLesson.FindOneAndUpdate(ctx, query, update, opts)

	var updatePost *lesson_domain.Lesson
	if err := res.Decode(&updatePost); err != nil {
		return nil, err
	}

	return updatePost, nil
}

func (l *lessonRepository) DeleteOne(ctx context.Context, lessonID string) error {
	collection := l.database.Collection(l.collectionLesson)
	objID, err := primitive.ObjectIDFromHex(lessonID)
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
		return errors.New(`the lesson is removed`)
	}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}
