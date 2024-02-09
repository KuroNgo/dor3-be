package repository_test

import (
	"clean-architecture/domain/request/user"
	mocks2 "clean-architecture/infrastructor/mocks"
	user_repository "clean-architecture/repository/user"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestCreate(t *testing.T) {

	var databaseHelper *mocks2.Database
	var collectionHelper *mocks2.Collection

	databaseHelper = &mocks2.Database{}
	collectionHelper = &mocks2.Collection{}

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
			On("InsertOne", mock.Anything, mock.AnythingOfType("*user_domain.User")).
			Return(mockUserID, nil).Once()

		databaseHelper.
			On("Collection", collectionName).
			Return(collectionHelper)

		ur := user_repository.NewUserRepository(databaseHelper, collectionName)

		err := ur.Create(context.Background(), mockUser)

		assert.NoError(t, err)

		collectionHelper.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		collectionHelper.
			On("InsertOne", mock.Anything, mock.AnythingOfType("*user_domain.User")).
			Return(mockEmptyUser, errors.New("unexpected")).Once()

		databaseHelper.
			On("Collection", collectionName).
			Return(collectionHelper)

		ur := user.NewUserRepository(databaseHelper, collectionName)

		// test trên hàm create
		err := ur.Create(context.Background(), mockEmptyUser)

		assert.Error(t, err)

		collectionHelper.AssertExpectations(t)
	})

}
