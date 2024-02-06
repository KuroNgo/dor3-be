package repository_test

import (
	"clean-architecture/domain"
	"clean-architecture/domain/request/user"
	"clean-architecture/infrastructor/mongo/mocks"
	"clean-architecture/repository"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestCreate(t *testing.T) {

	var databaseHelper *mocks.Database
	var collectionHelper *mocks.Collection

	databaseHelper = &mocks.Database{}
	collectionHelper = &mocks.Collection{}

	collectionName := domain.domain.CollectionUser

	mockUser := &user.User{
		ID:       primitive.NewObjectID(),
		FullName: "Test",
		Email:    "test@gmail.com",
		Password: "password",
	}

	mockEmptyUser := &user.User{}
	mockUserID := primitive.NewObjectID()

	t.Run("success", func(t *testing.T) {

		collectionHelper.
			On("InsertOne", mock.Anything, mock.AnythingOfType("*domain.User")).
			Return(mockUserID, nil).Once()

		databaseHelper.
			On("Collection", collectionName).
			Return(collectionHelper)

		ur := repository.NewUserRepository(databaseHelper, collectionName)

		err := ur.Create(context.Background(), mockUser)

		assert.NoError(t, err)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.
			On("InsertOne", mock.Anything, mock.AnythingOfType("*domain.User")).
			Return(mockEmptyUser, errors.New("unexpected")).Once()

		databaseHelper.
			On("Collection", collectionName).
			Return(collectionHelper)

		ur := repository.NewUserRepository(databaseHelper, collectionName)

		// test trên hàm create
		err := ur.Create(context.Background(), mockEmptyUser)

		assert.Error(t, err)

		collectionHelper.AssertExpectations(t)
	})

}
