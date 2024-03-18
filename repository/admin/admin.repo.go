package admin_repository

import (
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type adminRepository struct {
	database   mongo.Database
	collection string
}

func (a *adminRepository) GetByID(c context.Context, id string) (*admin_domain.Admin, error) {
	collection := a.database.Collection(a.collection)

	var admin admin_domain.Admin

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &admin, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&admin)
	return &admin, err
}

func (a *adminRepository) FetchMany(c context.Context) ([]admin_domain.Admin, error) {
	collection := a.database.Collection(a.collection)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var admin []admin_domain.Admin

	err = cursor.All(c, &admin)
	if admin == nil {
		return []admin_domain.Admin{}, err
	}

	return admin, err
}

func (a *adminRepository) GetByEmail(c context.Context, username string) (*admin_domain.Admin, error) {
	collection := a.database.Collection(a.collection)
	var admin admin_domain.Admin
	err := collection.FindOne(c, bson.M{"email": username}).Decode(&admin)

	return &admin, err
}

func (a *adminRepository) Login(c context.Context, request admin_domain.SignIn) (*admin_domain.Admin, error) {
	admin, err := a.GetByEmail(c, request.Email)

	// Kiểm tra xem mật khẩu đã nhập có đúng với mật khẩu đã hash trong cơ sở dữ liệu không
	if err = internal.VerifyPassword(admin.Password, request.Password); err != nil {
		return &admin_domain.Admin{}, errors.New("email or password not found! ")
	}
	return admin, nil
}

func (a *adminRepository) CreateOne(c context.Context, admin admin_domain.Admin) error {
	collection := a.database.Collection(a.collection)

	filter := bson.M{"email": admin.Email}
	count, err := collection.CountDocuments(c, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the email do not unique")
	}
	_, err = collection.InsertOne(c, admin)
	return err
}

func (a *adminRepository) UpdateOne(ctx context.Context, adminID string, admin admin_domain.Admin) error {
	collection := a.database.Collection(a.collection)
	doc, err := internal.ToDoc(admin)
	objID, err := primitive.ObjectIDFromHex(adminID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (a *adminRepository) DeleteOne(ctx context.Context, adminID string, admin admin_domain.Admin) error {
	collection := a.database.Collection(a.collection)
	objID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`the user is removed`)
	}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}

func (a *adminRepository) UpsertOne(c context.Context, email string, admin *admin_domain.Admin) error {
	collection := a.database.Collection(a.collection)
	doc, err := internal.ToDoc(admin)
	if err != nil {
		return err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collection.FindOneAndUpdate(c, query, update, opts)

	var updatedPost *user_domain.Response

	if err := res.Decode(&updatedPost); err != nil {
		return errors.New("no post with that Id exists")
	}

	return nil
}

func NewAdminRepository(db mongo.Database, collection string) admin_domain.IAdminRepository {
	return &adminRepository{
		database:   db,
		collection: collection,
	}
}
