package repository_test

import (
	"clean-architecture/domain/request/user"
	"clean-architecture/infrastructor/mongo/mocks"
	user_repository "clean-architecture/repository/user"
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

	collectionName := user_domain.CollectionUser

	mockUser := &user_domain.User{
		ID:       primitive.NewObjectID(),
		FullName: "Test",
		Email:    "test@gmail.com",
		Password: "password",
	}

	mockEmptyUser := &user_domain.User{}
	mockUserID := primitive.NewObjectID()

	t.Run("success", func(t *testing.T) {

		collectionHelper.
			On("FindOneAndUpdate", mock.Anything, mock.AnythingOfType("*user_domain.User")).
			Return(mockUserID, nil).Once()

		databaseHelper.
			On("Collection", collectionName).
			Return(collectionHelper)

		ur := user_repository.NewUserRepository(databaseHelper, collectionName)

		_, err := ur.UpsertUser(context.Background(), mockUser.Email, mockEmptyUser)

		assert.NoError(t, err)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.
			On("FindOneAndUpdate", mock.Anything, mock.AnythingOfType("*user_domain.User")).
			Return(mockEmptyUser, errors.New("unexpected")).Once()

		databaseHelper.
			On("Collection", collectionName).
			Return(collectionHelper)

		ur := user_repository.NewUserRepository(databaseHelper, collectionName)

		// test trên hàm upsert
		_, err := ur.UpsertUser(context.Background(), mockUser.Email, mockEmptyUser)

		assert.Error(t, err)

		collectionHelper.AssertExpectations(t)
	})

}
