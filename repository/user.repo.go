package repository

import (
	"clean-architecture/domain/request"
	"clean-architecture/infrastructor/mongo"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	database   mongo.Database
	collection string
}

func NewUserRepository(db mongo.Database, collection string) domain.IUserRepository {
	return &UserRepository{
		database:   db,
		collection: collection,
	}
}

func (u *UserRepository) Create(c context.Context, user *domain.User) error {
	collection := u.database.Collection(u.collection)
	_, err := collection.InsertOne(c, user)

	return err
}

func (u *UserRepository) CreateAsync(c context.Context, user *domain.User) <-chan error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) Fetch(c context.Context) ([]domain.User, error) {
	collection := u.database.Collection(u.collection)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var users []domain.User

	err = cursor.All(c, &users)
	if users == nil {
		return []domain.User{}, err
	}

	return users, err
}

func (u *UserRepository) Update(c context.Context, user *domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) Delete(userID primitive.ObjectID) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) GetByID(c context.Context, id primitive.ObjectID) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}
