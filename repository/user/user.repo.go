package user

import (
	"clean-architecture/domain/user"
	"clean-architecture/internal"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	database   *mongo.Database
	collection string
}

func (u *userRepository) UpdateVerifyForChangePassword(ctx context.Context, user *user_domain.User) (*mongo.UpdateResult, error) {
	collection := u.database.Collection(u.collection)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"verified":   user.Verified,
		"updated_at": user.UpdatedAt,
	}}}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *userRepository) UpdatePassword(ctx context.Context, user *user_domain.User) error {
	collection := u.database.Collection(u.collection)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"password":          user.Password,
		"verification_code": user.VerificationCode,
		"updated_at":        user.UpdatedAt,
	}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *mongo.Database, collection string) user_domain.IUserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

func (u *userRepository) Update(ctx context.Context, user *user_domain.User) error {
	collection := u.database.Collection(u.collection)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: user}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) CheckVerify(ctx context.Context, verificationCode string) bool {
	collection := u.database.Collection(u.collection)

	filter := bson.M{"verification_code": verificationCode}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil || count == 0 {
		return false
	}

	return true
}

func (u *userRepository) GetByVerificationCode(ctx context.Context, verificationCode string) (*user_domain.User, error) {
	collection := u.database.Collection(u.collection)

	filter := bson.M{"verification_code": verificationCode}

	var user user_domain.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) UpdateImage(c context.Context, userID string, imageURL string) error {
	collection := u.database.Collection(u.collection)
	doc, err := internal.ToDoc(imageURL)
	objID, err := primitive.ObjectIDFromHex(userID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collection.UpdateOne(c, filter, update)
	return err
}

func (u *userRepository) UpdateVerify(ctx context.Context, user *user_domain.User) (*mongo.UpdateResult, error) {
	collection := u.database.Collection(u.collection)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"verified":          user.Verified,
		"verification_code": user.VerificationCode,
		"updated_at":        user.UpdatedAt,
	}}}

	data, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *userRepository) Create(c context.Context, user *user_domain.User) error {
	collection := u.database.Collection(u.collection)

	filter := bson.M{"email": user.Email}
	count, err := collection.CountDocuments(c, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the email do not unique")
	}
	_, err = collection.InsertOne(c, &user)
	return err
}

func (u *userRepository) FetchMany(c context.Context) (user_domain.Response, error) {
	collection := u.database.Collection(u.collection)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return user_domain.Response{}, err
	}

	var users []user_domain.User

	err = cursor.All(c, &users)
	if users == nil {
		return user_domain.Response{}, err
	}

	return user_domain.Response{}, err
}

func (u *userRepository) DeleteOne(c context.Context, userID string) error {
	collection := u.database.Collection(u.collection)
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}
	_, err = collection.DeleteOne(c, filter)
	return err
}

func (u *userRepository) GetByEmail(c context.Context, email string) (*user_domain.User, error) {
	collection := u.database.Collection(u.collection)
	var user user_domain.User
	err := collection.FindOne(c, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (u *userRepository) Login(c context.Context, request user_domain.SignIn) (*user_domain.User, error) {
	user, err := u.GetByEmail(c, request.Email)

	// Kiểm tra xem mật khẩu đã nhập có đúng với mật khẩu đã hash trong cơ sở dữ liệu không
	if err = internal.VerifyPassword(user.Password, request.Password); err != nil {
		return &user_domain.User{}, errors.New("email or password not found! ")
	}
	return user, nil
}

func (u *userRepository) GetByID(c context.Context, id string) (*user_domain.User, error) {
	collection := u.database.Collection(u.collection)

	var user user_domain.User

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &user, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&user)
	return &user, err
}

func (u *userRepository) UpsertOne(c context.Context, email string, user *user_domain.User) (*user_domain.User, error) {
	collection := u.database.Collection(u.collection)

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"full_name":  user.FullName,
		"email":      user.Email,
		"password":   user.Password,
		"avatar_url": user.AvatarURL,
		"asset_id":   user.AssetID,
		"phone":      user.Phone,
		"provider":   user.Provider,
		"verified":   user.Verified,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"role":       user.Role,
	}}}
	res := collection.FindOneAndUpdate(c, filter, update, opts)

	var updatedUser *user_domain.User
	if err := res.Decode(&updatedUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if updatedUser != nil {
		return updatedUser, nil
	} else {
		return user, nil
	}
}

func (u *userRepository) UniqueVerificationCode(ctx context.Context, verificationCode string) bool {
	collection := u.database.Collection(u.collection)

	filter := bson.M{"verification_code": verificationCode}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil || count > 0 {
		return false
	}
	return true
}
