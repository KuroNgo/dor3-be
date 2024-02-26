package course_repository

import (
	course_domain "clean-architecture/domain/course"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type courseRepository struct {
	database   mongo.Database
	collection string
}

func NewCourseRepository(db mongo.Database, collection string) course_domain.ICourseRepository {
	return &courseRepository{
		database:   db,
		collection: collection,
	}
}

func (c *courseRepository) FetchByID(ctx context.Context, courseID string) (*course_domain.Course, error) {
	collection := c.database.Collection(c.collection)

	var course course_domain.Course

	idHex, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return &course, err
	}

	err = collection.
		FindOne(ctx, bson.M{"_id": idHex}).
		Decode(&course)
	return &course, err
}

func (c *courseRepository) FetchMany(ctx context.Context) ([]course_domain.Course, error) {
	collection := c.database.Collection(c.collection)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var course []course_domain.Course

	err = cursor.All(ctx, &course)
	if course == nil {
		return []course_domain.Course{}, err
	}

	return course, err
}

func (c *courseRepository) FetchToDeleteMany(ctx context.Context) (*[]course_domain.Course, error) {
	collection := c.database.Collection(c.collection)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var course *[]course_domain.Course
	err = cursor.All(ctx, course)
	if course == nil {
		return &[]course_domain.Course{}, err
	}

	return course, err
}

func (c *courseRepository) UpdateOne(ctx context.Context, courseID string, course course_domain.Course) error {
	collection := c.database.Collection(c.collection)
	doc, err := internal.ToDoc(course)
	objID, err := primitive.ObjectIDFromHex(courseID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (c *courseRepository) CreateOne(ctx context.Context, course *course_domain.Course) error {
	collection := c.database.Collection(c.collection)

	filter := bson.M{"name": course.Name}
	// check exists with CountDocuments
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the course name did exist")
	}

	_, err = collection.InsertOne(ctx, course)
	return err
}

func (c *courseRepository) UpsertOne(ctx context.Context, id string, course *course_domain.Course) (*course_domain.Response, error) {
	collection := c.database.Collection(c.collection)
	doc, err := internal.ToDoc(course)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collection.FindOneAndUpdate(ctx, query, update, opts)

	var updatedPost *course_domain.Response

	if err := res.Decode(&updatedPost); err != nil {
		return nil, errors.New("no post with that Id exists")
	}

	return updatedPost, nil
}

func (c *courseRepository) DeleteOne(ctx context.Context, courseID string) error {
	collection := c.database.Collection(c.collection)
	objID, err := primitive.ObjectIDFromHex(courseID)
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
	if count <= 0 {
		return errors.New(`the course is removed`)
	}
	_, err = collection.DeleteOne(ctx, filter)
	return err
}
