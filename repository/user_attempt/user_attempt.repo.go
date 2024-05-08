package user_attempt_repository

import (
	"clean-architecture/domain/user_attempt"
	"context"
	"errors"
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

	filterExam := bson.M{"_id": userProcess.ExamID}
	countExam, err := collectionExam.CountDocuments(ctx, filterExam)
	if err != nil {
		return err
	}
	if countExam == 0 {
		return errors.New("the exam do not exist")
	}

	filterQuiz := bson.M{"_id": userProcess.QuizID}
	countQuiz, err := collectionQuiz.CountDocuments(ctx, filterQuiz)
	if err != nil {
		return err
	}
	if countQuiz == 0 {
		return errors.New("the quiz do not exist")
	}

	filterExercise := bson.M{"_id": userProcess.ExerciseID}
	countExercise, err := collectionExercise.CountDocuments(ctx, filterExercise)
	if err != nil {
		return err
	}
	if countExercise == 0 {
		return errors.New("the quiz do not exist")
	}

	_, err = collectionUserAttempt.InsertOne(ctx, userProcess)
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
