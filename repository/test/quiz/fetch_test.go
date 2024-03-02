package quiz_test

import (
	quiz_domain "clean-architecture/domain/quiz"
	"clean-architecture/infrastructor/mongo/mocks"
	quiz_repository "clean-architecture/repository/quiz"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFetchOne(t *testing.T) {
	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}

	collectionName := quiz_domain.CollectionQuiz
	mockQuiz := quiz_domain.Quiz{
		ID:            primitive.NewObjectID(),
		Question:      "What is the capital of France?",
		Options:       []string{"Paris", "London", "Berlin", "Rome"},
		QuestionType:  "checkbox",
		CorrectAnswer: "Paris",
	}

	mockEmptyQuiz := &quiz_domain.Quiz{}
	mockQuizID := primitive.NewObjectID()

	t.Run("success", func(t *testing.T) {

		collectionHelper.On("FindOne", mock.Anything,
			mock.AnythingOfType("*quiz_domain.Quiz")).
			Return(mockQuizID, nil).
			Once()

		databaseHelper.
			On("Collection", collectionName).
			Return(collectionHelper)

		ur := quiz_repository.NewQuizRepository(databaseHelper, collectionName)

		quiz, err := ur.FetchByID(context.Background(), mockQuiz.ID.Hex())

		assert.NoError(t, err)

		assert.NotNil(t, quiz)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.On("FindOne", mock.Anything,
			mock.AnythingOfType("*quiz_domain.Quiz")).
			Return(mockEmptyQuiz, errors.New("unexpected")).
			Once()
		databaseHelper.
			On("Collection", collectionName).
			Return(collectionHelper)

		ur := quiz_repository.NewQuizRepository(databaseHelper, collectionName)

		quiz, err := ur.FetchByID(context.Background(), mockEmptyQuiz.ID.Hex())

		assert.Error(t, err)

		assert.Nil(t, quiz)

		collectionHelper.AssertExpectations(t)
	})

}
