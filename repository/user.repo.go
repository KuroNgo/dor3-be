package repository

import (
	"clean-architecture/domain/request/user"
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

func NewUserRepository(db mongo.Database, collection string) user_domain.IUserRepository {
	return &UserRepository{
		database:   db,
		collection: collection,
	}
}

// Create interacted with user in domain to database
func (u *UserRepository) Create(c context.Context, user *user_domain.User) error {
	collection := u.database.Collection(u.collection)
	_, err := collection.InsertOne(c, user)

	return err
}

func (u *UserRepository) CreateAsync(c context.Context, user *user_domain.User) <-chan error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) Fetch(c context.Context) ([]user_domain.User, error) {
	collection := u.database.Collection(u.collection)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var users []user_domain.User

	err = cursor.All(c, &users)
	if users == nil {
		return []user_domain.User{}, err
	}

	return users, err
}

func (u *UserRepository) Update(c context.Context, userID primitive.ObjectID, updatedUser interface{}) error {
	collection := u.database.Collection(u.collection)
	objID, err := primitive.ObjectIDFromHex(userID.Hex())
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}
	update := bson.M{
		"$set": updatedUser,
	}

	_, err = collection.UpdateOne(c, filter, update)
	return err
}

func (u *UserRepository) Delete(c context.Context, userID primitive.ObjectID) error {
	collection := u.database.Collection(u.collection)
	objID, err := primitive.ObjectIDFromHex(userID.Hex())
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}
	_, err = collection.DeleteOne(c, filter)
	return err
}

func (u *UserRepository) GetByEmail(c context.Context, email string) (user_domain.User, error) {
	collection := u.database.Collection(u.collection)
	var user user_domain.User
	err := collection.FindOne(c, bson.M{"email": email}).Decode(&user)
	return user, err
}

func (u *UserRepository) GetByUsername(c context.Context, username string) (user_domain.User, error) {
	collection := u.database.Collection(u.collection)
	var user user_domain.User
	err := collection.FindOne(c, bson.M{"email": username}).Decode(&user)
	return user, err
}

func (u *UserRepository) GetByID(c context.Context, id primitive.ObjectID) (user_domain.User, error) {
	collection := u.database.Collection(u.collection)

	var user user_domain.User

	idHex, err := primitive.ObjectIDFromHex(id.Hex())
	if err != nil {
		return user, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&user)
	return user, err
}
