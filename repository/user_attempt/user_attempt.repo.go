package user_attempt_repository

import (
	"clean-architecture/domain/user_attempt"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userAttemptRepository struct {
	database              *mongo.Database
	collectionExam        string
	collectionExercise    string
	collectionQuiz        string
	collectionUserAttempt string
}

func NewUserAttemptRepository(db *mongo.Database, collectionUserAttempt string, collectionExam string, collectionQuiz string, collectionExercise string) user_attempt_domain.IUserProcessRepository {
	return &userAttemptRepository{
		database:              db,
		collectionUserAttempt: collectionUserAttempt,
		collectionExam:        collectionExam,
		collectionExercise:    collectionExercise,
		collectionQuiz:        collectionQuiz,
	}
}

func (u *userAttemptRepository) FetchManyByUserID(ctx context.Context, userID string) (user_attempt_domain.Response, error) {
	collectionExam := u.database.Collection(u.collectionUserAttempt)

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user_attempt_domain.Response{}, err
	}

	filter := bson.M{"user_id": idUser}

	var userAttempt user_attempt_domain.UserProcess
	err = collectionExam.FindOne(ctx, filter).Decode(userAttempt)
	if err != nil {
		return user_attempt_domain.Response{}, err
	}

	userAttemptRes := user_attempt_domain.Response{
		UserProcess: userAttempt,
	}
	return userAttemptRes, nil

}

func (u *userAttemptRepository) CreateOneByUserID(ctx context.Context, userProcess user_attempt_domain.UserProcess) error {
	collectionUserAttempt := u.database.Collection(u.collectionUserAttempt)
	collectionExam := u.database.Collection(u.collectionExam)
	collectionQuiz := u.database.Collection(u.collectionQuiz)
	collectionExercise := u.database.Collection(u.collectionExercise)

	// Hàm kiểm tra sự tồn tại của tài liệu
	checkExistence := func(collection *mongo.Collection, id primitive.ObjectID, idName string) error {
		if id == primitive.NilObjectID {
			return nil
		}
		filter := bson.M{"_id": id}
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("the %s does not exist", idName)
		}
		return nil
	}

	// Kiểm tra sự tồn tại của Exam
	if err := checkExistence(collectionExam, userProcess.ExamID, "exam"); err != nil {
		return err
	}

	// Kiểm tra sự tồn tại của Quiz
	if err := checkExistence(collectionQuiz, userProcess.QuizID, "quiz"); err != nil {
		return err
	}

	// Kiểm tra sự tồn tại của Exercise
	if err := checkExistence(collectionExercise, userProcess.ExerciseID, "exercise"); err != nil {
		return err
	}

	// Chèn userProcess vào collectionUserAttempt
	_, err := collectionUserAttempt.InsertOne(ctx, userProcess)
	return err
}

func (u *userAttemptRepository) UpdateAttemptByUserID(ctx context.Context, user user_attempt_domain.UserProcess) error {
	collectionUserAttempt := u.database.Collection(u.collectionUserAttempt)

	filter := bson.D{{Key: "user_id", Value: user.UserID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"score":          user.Score,
		"process_status": user.ProcessStatus,
		"completed_date": user.CompletedDate,
		"updated_at":     user.UpdatedAt,
	}}}

	_, err := collectionUserAttempt.UpdateOne(ctx, filter, &update)
	return err
}

func (u *userAttemptRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	collectionUserAttempt := u.database.Collection(u.collectionUserAttempt)

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"user_id": idUser}

	_, err = collectionUserAttempt.DeleteOne(ctx, &filter)
	return err
}
