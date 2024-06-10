package user

import (
	"clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	"clean-architecture/internal"
	"clean-architecture/internal/cache"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type userRepository struct {
	database             *mongo.Database
	collectionUser       string
	collectionUserDetail string
}

func NewUserRepository(db *mongo.Database, collectionUser string, collectionUserDetail string) user_domain.IUserRepository {
	return &userRepository{
		database:             db,
		collectionUser:       collectionUser,
		collectionUserDetail: collectionUserDetail,
	}
}

var (
	userCache  = cache.New[string, *user_domain.User]()
	usersCache = cache.NewTTL[string, user_domain.Response]()

	wg sync.WaitGroup
	mu sync.Mutex
)

func (u *userRepository) UpdateVerifyForChangePassword(ctx context.Context, user *user_domain.User) (*mongo.UpdateResult, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"verified":   user.Verified,
		"updated_at": user.UpdatedAt,
	}}}

	data, err := collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *userRepository) UpdatePassword(ctx context.Context, user *user_domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"password":          user.Password,
		"verification_code": user.VerificationCode,
		"updated_at":        user.UpdatedAt,
	}}}

	filterUnique := bson.M{"email": user.Email}
	count, err := collectionUser.CountDocuments(ctx, filterUnique)
	if count > 0 {
		return errors.New("the email must be unique")
	}
	_, err = collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) Update(ctx context.Context, user *user_domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"full_name":  user.FullName,
		"phone":      user.Phone,
		"updated_at": user.UpdatedAt,
	}}}

	_, err := collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		userCache.Clear()
	}()
	wg.Wait()
	return nil
}

func (u *userRepository) CheckVerify(ctx context.Context, verificationCode string) bool {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"verification_code": verificationCode}
	count, err := collectionUser.CountDocuments(ctx, filter)
	if err != nil || count == 0 {
		return false
	}

	return true
}

func (u *userRepository) GetByVerificationCode(ctx context.Context, verificationCode string) (*user_domain.User, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"verification_code": verificationCode}

	var user user_domain.User
	err := collectionUser.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) UpdateImage(c context.Context, userID string, imageURL string) error {
	collectionUser := u.database.Collection(u.collectionUser)
	doc, err := internal.ToDoc(imageURL)
	objID, err := primitive.ObjectIDFromHex(userID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err = collectionUser.UpdateOne(c, filter, update)

	wg.Add(1)
	go func() {
		defer wg.Done()
		userCache.Clear()
	}()
	wg.Wait()

	return err
}

func (u *userRepository) UpdateVerify(ctx context.Context, user *user_domain.User) (*mongo.UpdateResult, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"verified":          user.Verified,
		"verification_code": user.VerificationCode,
		"updated_at":        user.UpdatedAt,
	}}}

	data, err := collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *userRepository) Create(c context.Context, user *user_domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"email": user.Email}
	count, err := collectionUser.CountDocuments(c, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the email do not unique")
	}
	_, err = collectionUser.InsertOne(c, &user)

	wg.Add(1)
	go func() {
		defer wg.Done()
		usersCache.Clear()
	}()
	wg.Wait()

	return err
}

func (u *userRepository) FetchMany(c context.Context) (user_domain.Response, error) {
	usersCh := make(chan user_domain.Response)
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := usersCache.Get("users")
		if found {
			usersCh <- data
			return
		}
	}()

	go func() {
		defer close(usersCh)
		wg.Wait()
	}()

	userData := <-usersCh
	if !internal.IsZeroValue(userData) {
		return userData, nil
	}

	collectionUser := u.database.Collection(u.collectionUser)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collectionUser.Find(c, bson.D{}, opts)

	if err != nil {
		return user_domain.Response{}, err
	}

	var users []user_domain.User

	err = cursor.All(c, &users)
	if users == nil {
		return user_domain.Response{}, err
	}

	var statisticsCh = make(chan user_domain.Statistics)
	go func() {
		defer close(statisticsCh)
		statistics, _ := u.Statistics(c)
		statisticsCh <- statistics
	}()

	statistics := <-statisticsCh
	response := user_domain.Response{
		User:       users,
		Statistics: statistics,
	}

	usersCache.Set("users", response, 5*time.Minute)
	return response, err
}

func (u *userRepository) DeleteOne(c context.Context, userID string) error {
	collectionUser := u.database.Collection(u.collectionUser)
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}
	_, err = collectionUser.DeleteOne(c, filter)

	wg.Add(2)
	go func() {
		defer wg.Done()
		userCache.Clear()
	}()

	go func() {
		defer wg.Done()
		usersCache.Clear()
	}()
	wg.Wait()

	return err
}

func (u *userRepository) GetByEmail(c context.Context, email string) (*user_domain.User, error) {
	collectionUser := u.database.Collection(u.collectionUser)
	var user user_domain.User
	err := collectionUser.FindOne(c, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (u *userRepository) Login(c context.Context, request user_domain.SignIn) (*user_domain.User, error) {
	userCh := make(chan *user_domain.User)
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := userCache.Get(request.Email + request.Password)
		if found {
			userCh <- data
			return
		}
	}()

	go func() {
		defer close(userCh)
		wg.Wait()
	}()

	userData := <-userCh
	if !internal.IsZeroValue(userData) {
		return userData, nil
	}

	user, err := u.GetByEmail(c, request.Email)

	// Kiểm tra xem mật khẩu đã nhập có đúng với mật khẩu đã hash trong cơ sở dữ liệu không
	if err = internal.VerifyPassword(user.Password, request.Password); err != nil {
		return &user_domain.User{}, errors.New("email or password not found! ")
	}

	userCache.Set(request.Email+request.Password, user)
	return user, nil
}

func (u *userRepository) GetByID(c context.Context, id string) (*user_domain.User, error) {
	userCh := make(chan *user_domain.User)
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := userCache.Get(id)
		if found {
			userCh <- data
			return
		}
	}()

	go func() {
		defer close(userCh)
		wg.Wait()
	}()

	userData := <-userCh
	if !internal.IsZeroValue(userData) {
		return userData, nil
	}

	collectionUser := u.database.Collection(u.collectionUser)

	var user user_domain.User

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &user, err
	}

	err = collectionUser.FindOne(c, bson.M{"_id": idHex}).Decode(&user)
	if err != nil {
		return nil, err
	}
	userCache.Set(id, &user)
	return &user, nil
}

func (u *userRepository) UpsertOne(c context.Context, email string, user *user_domain.User) (*user_domain.User, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"full_name":  user.FullName,
		"email":      user.Email,
		"avatar_url": user.AvatarURL,
		"phone":      user.Phone,
		"provider":   user.Provider,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"role":       user.Role,
	}}}
	res := collectionUser.FindOneAndUpdate(c, filter, update, opts)

	var updatedUser *user_domain.User
	if err := res.Decode(&updatedUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		userCache.Clear()
	}()
	wg.Wait()

	return updatedUser, nil
}

func (u *userRepository) UniqueVerificationCode(ctx context.Context, verificationCode string) bool {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"verification_code": verificationCode}
	count, err := collectionUser.CountDocuments(ctx, filter)
	if err != nil || count > 0 {
		return false
	}
	return true
}

func (u *userRepository) Statistics(ctx context.Context) (user_domain.Statistics, error) {
	collectionUser := u.database.Collection(u.collectionUser)
	collectionUserDetail := u.database.Collection(u.collectionUserDetail)

	// Đếm tổng số lượng tài liệu trong collection
	count, err := collectionUser.CountDocuments(ctx, bson.D{})
	if err != nil {
		return user_domain.Statistics{}, err
	}

	cursor, err := collectionUserDetail.Find(ctx, bson.D{})
	if err != nil {
		return user_domain.Statistics{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var (
		countOutside int64 = 0
		countInside  int64 = 0
		countStudent int64 = 0
	)
	for cursor.Next(ctx) {
		var user user_detail_domain.UserDetail
		if err := cursor.Decode(&user); err != nil {
			return user_domain.Statistics{}, err
		}

		if user.Specialize == "inside" {
			countInside++
		}
		if user.Specialize == "outside" {
			countOutside++
		}
		if user.Specialize == "student" {
			countStudent++
		}
	}

	statistics := user_domain.Statistics{
		CountInside:  countInside,
		CountOutside: countOutside,
		CountStudent: countStudent,
		Total:        count,
	}

	return statistics, nil
}
