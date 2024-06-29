package course_repository

import (
	course_domain "clean-architecture/domain/course"
	"clean-architecture/infrastructor/mongo"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

const (
	course        = "course"
	courseProcess = "courseProcess"
	lesson        = "lesson"
	unit          = "unit"
	vocabulary    = "vocabulary"
)

func TestUpdateOneInAdmin(t *testing.T) {
	// Setup test MongoDB database
	client, database := mongo.SetupTestDatabase(t)
	defer mongo.TearDownTestDatabase(client, t)

	// Access the collection for courses
	collectionCourse := database.Collection(course)

	// Retrieve the first course from the database
	var firstCourse course_domain.Course
	err := collectionCourse.FindOne(
		context.Background(),
		bson.M{},
		options.FindOne().SetSort(bson.D{{Key: "_id", Value: 1}}),
	).Decode(&firstCourse)
	assert.NoError(t, err, "Error fetching first course")

	// Prepare mock data for testing
	mockCourse := &course_domain.Course{
		Id:          firstCourse.Id,
		Name:        "Test Course",
		Description: "This is a test course description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		WhoUpdated:  "tester",
	}

	// Test case for successful update
	t.Run("success", func(t *testing.T) {
		ur := NewCourseRepository(database, course, courseProcess, lesson, unit, vocabulary)
		_, err := ur.UpdateOneInAdmin(context.Background(), mockCourse)
		assert.NoError(t, err, "Error updating course")
	})

}

func TestCreateOneInAdmin(t *testing.T) {
	client, database := mongo.SetupTestDatabase(t)
	defer mongo.TearDownTestDatabase(client, t)

	mockCourse := &course_domain.Course{
		Id:          primitive.NewObjectID(),
		Name:        "Test Course",
		Description: "This is a test course description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		WhoUpdated:  "tester",
	}
	mockEmptyCourse := &course_domain.Course{}

	t.Run("success", func(t *testing.T) {
		ur := NewCourseRepository(database, course, courseProcess, lesson, unit, vocabulary)

		_ = ur.CreateOneInAdmin(context.Background(), mockCourse)
	})

	t.Run("error", func(t *testing.T) {
		ur := NewCourseRepository(database, course, courseProcess, lesson, unit, vocabulary)

		// Trying to insert an empty user, expecting an error
		err := ur.CreateOneInAdmin(context.Background(), mockEmptyCourse)
		assert.Error(t, err)
	})
}
