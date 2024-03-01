package lesson

import (
	lesson_domain "clean-architecture/domain/lesson"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo_drivermongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type lessonRepository struct {
	database   mongo.Database
	collection string
}

func NewLessonRepository(db mongo.Database, collection string) lesson_domain.ILessonRepository {
	return &lessonRepository{
		database:   db,
		collection: collection,
	}
}

func (l *lessonRepository) FetchByID(ctx context.Context, lessonID string) (*lesson_domain.Lesson, error) {
	collection := l.database.Collection(l.collection)

	var lesson lesson_domain.Lesson

	idHex, err := primitive.ObjectIDFromHex(lessonID)
	if err != nil {
		return &lesson, err
	}

	err = collection.
		FindOne(ctx, bson.M{"_id": idHex}).
		Decode(&lesson)
	return &lesson, err
}

func (l *lessonRepository) FetchMany(ctx context.Context) ([]lesson_domain.Lesson, error) {
	collection := l.database.Collection(l.collection)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var lesson []lesson_domain.Lesson
	err = cursor.All(ctx, &lesson)
	if lesson == nil {
		return []lesson_domain.Lesson{}, err
	}
	return lesson, err
}

func (l *lessonRepository) FetchToDeleteMany(ctx context.Context) (*[]lesson_domain.Lesson, error) {
	collection := l.database.Collection(l.collection)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var lesson *[]lesson_domain.Lesson
	err = cursor.All(ctx, lesson)
	if lesson == nil {
		return &[]lesson_domain.Lesson{}, err
	}
	return lesson, err
}

func (l *lessonRepository) UpdateOne(ctx context.Context, lessonID string, lesson lesson_domain.Lesson) error {
	collection := l.database.Collection(l.collection)
	doc, err := internal.ToDoc(lesson)
	objID, err := primitive.ObjectIDFromHex(lessonID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (l *lessonRepository) CreateOne(ctx context.Context, lesson *lesson_domain.Lesson) error {
	collection := l.database.Collection(l.collection)

	filter := bson.M{"name": lesson.Name}
	// check exists with CountDocuments
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the lesson name did exist")
	}

	// this sentence just only processed one collection
	//matchStage := bson.D{{Key: "$match", Value: bson.D{{"course_id", lesson.CourseID}}}}
	//groupStage := bson.D{{"$group", bson.D{{"_id", "$course_id"}}}}

	// this sentence can be processed to reference courseID
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "course"}, {"localField", "course"}, {"foreignField", "_id"}, {"as", "course"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$course"}, {"preserveNullAndEmptyArrays", false}}}}
	showInfoCursor, err := collection.Aggregate(ctx, mongo_drivermongo.Pipeline{lookupStage, unwindStage})
	if err != nil {
		return err
	}
	var showLoaded *[]bson.M
	if err = showInfoCursor.All(ctx, showLoaded); err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, lesson)
	return nil
}

func (l *lessonRepository) UpsertOne(ctx context.Context, id string, lesson *lesson_domain.Lesson) (*lesson_domain.Lesson, error) {
	collection := l.database.Collection(l.collection)
	doc, err := internal.ToDoc(lesson)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collection.FindOneAndUpdate(ctx, query, update, opts)

	var updatePost *lesson_domain.Lesson
	if err := res.Decode(&updatePost); err != nil {
		return nil, err
	}

	return updatePost, nil
}

func (l *lessonRepository) DeleteOne(ctx context.Context, lessonID string) error {
	collection := l.database.Collection(l.collection)
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
	if count < 0 {
		return errors.New(`the lesson is removed`)
	}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}
