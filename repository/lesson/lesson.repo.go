package lesson_repository

import (
	lesson_domain "clean-architecture/domain/lesson"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type lessonRepository struct {
	database         *mongo.Database
	collectionLesson string
	collectionCourse string
	collectionUnit   string
}

func NewLessonRepository(db *mongo.Database, collectionLesson string, collectionCourse string, collectionUnit string) lesson_domain.ILessonRepository {
	return &lessonRepository{
		database:         db,
		collectionLesson: collectionLesson,
		collectionCourse: collectionCourse,
		collectionUnit:   collectionUnit,
	}
}

func (l *lessonRepository) FetchByIdCourse(ctx context.Context, idCourse string) (lesson_domain.Response, error) {
	collectionLesson := l.database.Collection(l.collectionLesson)

	idCourse2, err := primitive.ObjectIDFromHex(idCourse)
	if err != nil {
		return lesson_domain.Response{}, err
	}

	filter := bson.M{"course_id": idCourse2}

	cursor, err := collectionLesson.Find(ctx, filter)
	if err != nil {
		return lesson_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var lessons []lesson_domain.Lesson

	for cursor.Next(ctx) {
		var lesson lesson_domain.Lesson
		if err = cursor.Decode(&lesson); err != nil {
			return lesson_domain.Response{}, err
		}

		// Gắn CourseID vào bài học
		lesson.CourseID = idCourse2

		lessons = append(lessons, lesson)
	}

	response := lesson_domain.Response{
		Lesson: lessons,
	}
	return response, nil
}

func (l *lessonRepository) FindCourseIDByCourseName(ctx context.Context, courseName string) (primitive.ObjectID, error) {
	collectionCourse := l.database.Collection(l.collectionCourse)

	filter := bson.M{"name": courseName}
	var data struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	err := collectionCourse.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return data.Id, nil
}

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
	return nil
}

func (l *lessonRepository) FetchMany(ctx context.Context) (lesson_domain.Response, error) {
	collectionLesson := l.database.Collection(l.collectionLesson)

	cursor, err := collectionLesson.Find(ctx, bson.D{})
	if err != nil {
		return lesson_domain.Response{}, err
	}

	var lessons []lesson_domain.Lesson
	for cursor.Next(ctx) {
		var lesson lesson_domain.Lesson
		if err = cursor.Decode(&lesson); err != nil {
			return lesson_domain.Response{}, err
		}

		// Thêm lesson vào slice lessons
		lessons = append(lessons, lesson)
	}
	lessonRes := lesson_domain.Response{
		Lesson: lessons,
	}

	return lessonRes, err
}

func (l *lessonRepository) UpdateOne(ctx context.Context, lesson *lesson_domain.Lesson) (*mongo.UpdateResult, error) {
	collection := l.database.Collection(l.collectionLesson)

	filter := bson.M{"_id": lesson.ID}

	update := bson.M{
		"$set": bson.M{
			"name":    lesson.Name,
			"content": lesson.Content,
			"image":   lesson.ImageURL,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	return data, err
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

func (l *lessonRepository) CreateOneByNameCourse(ctx context.Context, lesson *lesson_domain.Lesson) error {
	collectionLesson := l.database.Collection(l.collectionLesson)

	filter := bson.M{"name": lesson.Name}
	// check exists with CountDocuments
	count, err := collectionLesson.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the lesson name did exist")
	}

	_, err = collectionLesson.InsertOne(ctx, lesson)
	return nil
}

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
	return err
}
