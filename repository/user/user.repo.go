package user_repository

import (
	"clean-architecture/domain/request/user"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	database   mongo.Database
	collection string
}

func NewUserRepository(db mongo.Database, collection string) user_domain.IUserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

// Create interacted with user in domain to database
func (u *userRepository) Create(c context.Context, user *user_domain.User) error {
	collection := u.database.Collection(u.collection)
	_, err := collection.InsertOne(c, user)

	return err
}

func (u *userRepository) CreateAsync(c context.Context, user *user_domain.User) <-chan error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) Fetch(c context.Context) ([]user_domain.User, error) {
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

func (u *userRepository) Update(c context.Context, userID primitive.ObjectID, updatedUser interface{}) error {
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

func (u *userRepository) Delete(c context.Context, userID primitive.ObjectID) error {
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

func (u *userRepository) GetByEmail(c context.Context, email string) (user_domain.User, error) {
	collection := u.database.Collection(u.collection)
	var user user_domain.User
	err := collection.FindOne(c, bson.M{"email": email}).Decode(&user)
	return user, err
}

func (u *userRepository) GetByUsername(c context.Context, username string) (user_domain.User, error) {
	collection := u.database.Collection(u.collection)
	var user user_domain.User
	err := collection.FindOne(c, bson.M{"email": username}).Decode(&user)
	return user, err
}

func (u *userRepository) GetByID(c context.Context, id primitive.ObjectID) (user_domain.User, error) {
	collection := u.database.Collection(u.collection)

	var user user_domain.User

	idHex, err := primitive.ObjectIDFromHex(id.Hex())
	if err != nil {
		return user, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&user)
	return user, err
}

func (u *userRepository) UpsertUser(c context.Context, email string, user *user_domain.User) (*user_domain.Response, error) {
	collection := u.database.Collection(u.collection)
	doc, err := internal.ToDoc(user)
	if err != nil {
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collection.FindOneAndUpdate(c, query, update, opts)

	var updatedPost *user_domain.Response

	if err := res.Decode(&updatedPost); err != nil {
		return nil, errors.New("no post with that Id exists")
	}

	return updatedPost, nil
}