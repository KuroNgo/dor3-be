package quiz_test

import (
	quiz_domain "clean-architecture/domain/quiz"
	"clean-architecture/infrastructor/mongo/mocks"
	quiz_repository "clean-architecture/repository"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreate(t *testing.T) {
	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}

	collectionName := quiz_domain.CollectionQuiz
	mockQuiz := &quiz_domain.Input{
		Question:      "What is the capital of France?",
		Options:       []string{"Paris", "London", "Berlin", "Rome"},
		QuestionType:  "checkbox",
		CorrectAnswer: "Paris",
	}

	mockEmptyQuiz := &quiz_domain.Input{}
	mockQuizQuestion := "What is the capital of France?"

	t.Run("success", func(t *testing.T) {
		collectionHelper.On("InsertOne", mock.Anything,
			mock.AnythingOfType("*quiz_domain.Input")).
			Return(mockQuizQuestion, nil).
			Once()
		databaseHelper.
			On("Collection", collectionName).
			Return(collectionHelper)

		ur := quiz_repository.NewQuizRepository(databaseHelper, collectionName)

		err := ur.CreateOne(context.Background(), mockQuiz)

		assert.NoError(t, err)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.On("InsertOne", mock.Anything,
			mock.AnythingOfType("*quiz_domain.Input")).
			Return(mockEmptyQuiz, errors.New("Unexpected")).
			Once()
		databaseHelper.
			On("Collection", collectionName).
			Return(collectionHelper)

		ur := quiz_repository.NewQuizRepository(databaseHelper, collectionName)

		err := ur.CreateOne(context.Background(), mockEmptyQuiz)

		assert.Error(t, err)

		collectionHelper.AssertExpectations(t)
	})
}
