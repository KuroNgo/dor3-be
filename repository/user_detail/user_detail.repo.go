package user_detail_repository

import (
	user_detail_domain "clean-architecture/domain/user_detail"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type userDetailRepository struct {
	database   *mongo.Database
	collection string
}

func NewUserDetailRepository(db *mongo.Database, collection string) user_detail_domain.IUserDetailRepository {
	return &userDetailRepository{
		database:   db,
		collection: collection,
	}
}

func (u *userDetailRepository) FetchByUserID(ctx context.Context, userid string) (user_detail_domain.UserDetail, error) {
	collection := u.database.Collection(u.collection)

	idUser, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return user_detail_domain.UserDetail{}, err
	}

	filter := bson.M{"user_id": idUser}
	var user user_detail_domain.UserDetail
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user_detail_domain.UserDetail{}, err
	}

	return user, nil
}

func (u *userDetailRepository) Create(ctx context.Context, user user_detail_domain.UserDetail) error {
	collection := u.database.Collection(u.collection)

	filter := bson.M{"name": user.UserID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the course name did exist")
	}

	_, err = collection.InsertOne(ctx, user)
	return err
}

func (u *userDetailRepository) Update(ctx context.Context, user *user_detail_domain.UserDetail) (*mongo.UpdateResult, error) {
	collection := u.database.Collection(u.collection)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.M{
		"$set": bson.M{
			"specialize": user.Specialize,
		},
	}

	var mu sync.Mutex // Mutex để bảo vệ courses

	mu.Lock()
	data, err := collection.UpdateOne(ctx, filter, &update)
	mu.Unlock()

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (u *userDetailRepository) DeleteOne(ctx context.Context, userID string) error {
	collection := u.database.Collection(u.collection)

	filter := bson.M{"_id": userID}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
