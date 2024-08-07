package admin_repository

import (
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal"
	"clean-architecture/internal/cache/memory"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type adminRepository struct {
	database        *mongo.Database
	collectionAdmin string
	collectionUser  string
}

func NewAdminRepository(db *mongo.Database, collectionAdmin string, collectionUser string) admin_domain.IAdminRepository {
	return &adminRepository{
		database:        db,
		collectionAdmin: collectionAdmin,
		collectionUser:  collectionUser,
	}
}

var (
	adminsCache = memory.NewTTL[string, admin_domain.Response]()
	wg          sync.WaitGroup
)

func (a *adminRepository) GetByID(c context.Context, id string) (*admin_domain.Admin, error) {
	collection := a.database.Collection(a.collectionAdmin)

	var admin admin_domain.Admin
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &admin, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&admin)
	return &admin, err
}

func (a *adminRepository) FetchMany(c context.Context) (admin_domain.Response, error) {
	errCh := make(chan error)
	adminsCh := make(chan admin_domain.Response)
	wg.Add(1)
	go func() {
		data, found := adminsCache.Get("admins")
		if found {
			adminsCh <- data
		}
	}()
	wg.Wait()

	adminData := <-adminsCh
	if !internal.IsZeroValue(adminData) {
		return adminData, nil
	}

	collection := a.database.Collection(a.collectionAdmin)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)
	if err != nil {
		return admin_domain.Response{}, err
	}

	var admins []admin_domain.Admin

	wg.Add(1)
	go func() {
		defer wg.Done()
		for cursor.Next(c) {
			var admin admin_domain.Admin
			if err = cursor.Decode(&admin); err != nil {
				errCh <- err
				return
			}

			admin.CreatedAt = admin.CreatedAt.Add(7 + time.Hour)
			admin.UpdatedAt = admin.UpdatedAt.Add(7 + time.Hour)

			// Thêm lesson vào slice lessons
			admins = append(admins, admin)
		}
	}()

	wg.Wait()

	statisticsCh := make(chan admin_domain.Statistics)
	go func() {
		statistics, _ := a.Statistics(c)
		statisticsCh <- statistics
	}()

	statistics := <-statisticsCh

	adminRes := admin_domain.Response{
		Admin:      admins,
		Statistics: statistics,
	}

	adminsCache.Set("admins", adminRes, 5*time.Minute)
	select {
	case err = <-errCh:
		return admin_domain.Response{}, err
	default:
		return adminRes, err
	}
}

func (a *adminRepository) GetByEmail(c context.Context, username string) (*admin_domain.Admin, error) {
	collection := a.database.Collection(a.collectionAdmin)
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
	collectionAdmin := a.database.Collection(a.collectionAdmin)
	collectionUser := a.database.Collection(a.collectionUser)

	filter := bson.M{"email": admin.Email}
	count, err := collectionAdmin.CountDocuments(c, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the email do not unique")
	}

	count, err = collectionUser.CountDocuments(c, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the email do not the same email in user")
	}

	_, err = collectionAdmin.InsertOne(c, admin)

	wg.Add(1)
	go func() {
		defer wg.Done()
		adminsCache.Clear()
	}()
	wg.Wait()

	return err
}

func (a *adminRepository) UpdateOne(ctx context.Context, admin *admin_domain.Admin) (*mongo.UpdateResult, error) {
	collection := a.database.Collection(a.collectionAdmin)

	filter := bson.M{"_id": admin.Id}
	update := bson.M{
		"$set": bson.M{
			"full_name":  admin.FullName,
			"address":    admin.Address,
			"phone":      admin.Phone,
			"updated_at": admin.UpdatedAt,
			"avatar":     admin.AvatarURL,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		adminsCache.Clear()
	}()
	wg.Wait()

	return data, err
}

func (a *adminRepository) ChangeEmail(ctx context.Context, admin *admin_domain.Admin) (*mongo.UpdateResult, error) {
	collection := a.database.Collection(a.collectionAdmin)

	filter := bson.M{"_id": admin.Id}
	update := bson.M{
		"$set": bson.M{
			"email":      admin.Phone,
			"updated_at": admin.UpdatedAt,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (a *adminRepository) DeleteOne(ctx context.Context, adminID string) error {
	collection := a.database.Collection(a.collectionAdmin)
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
		return errors.New(`the admin is removed`)
	}

	countElement, err := collection.CountDocuments(ctx, bson.M{})
	if countElement <= 1 {
		return errors.New(`the admin can not delete`)
	}
	_, err = collection.DeleteOne(ctx, filter)

	wg.Add(1)
	go func() {
		defer wg.Done()
		adminsCache.Clear()
	}()
	wg.Wait()
	return err
}

func (a *adminRepository) UpsertOne(c context.Context, email string, admin *admin_domain.Admin) error {
	collectionAdmin := a.database.Collection(a.collectionAdmin)
	collectionUser := a.database.Collection(a.collectionUser)

	doc, err := internal.ToDoc(admin)
	if err != nil {
		return err
	}

	count, err := collectionUser.CountDocuments(c, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the email do not the same email in user")
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collectionAdmin.FindOneAndUpdate(c, query, update, opts)

	var updatedPost *user_domain.Response

	if err = res.Decode(&updatedPost); err != nil {
		return errors.New("no post with that Id exists")
	}

	return nil
}

func (a *adminRepository) Statistics(ctx context.Context) (admin_domain.Statistics, error) {
	collectionAdmin := a.database.Collection(a.collectionAdmin)

	count, err := collectionAdmin.CountDocuments(ctx, bson.M{})
	if err != nil {
		return admin_domain.Statistics{}, err
	}

	statistics := admin_domain.Statistics{
		Total: count,
	}

	return statistics, nil
}
