package user_attempt_repository

import (
	"clean-architecture/domain/user_process/exam_management"
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

func NewUserAttemptRepository(db *mongo.Database, collectionUserAttempt string, collectionExam string, collectionQuiz string, collectionExercise string) exam_management.IUserProcessRepository {
	return &userAttemptRepository{
		database:              db,
		collectionUserAttempt: collectionUserAttempt,
		collectionExam:        collectionExam,
		collectionExercise:    collectionExercise,
		collectionQuiz:        collectionQuiz,
	}
}

func (u *userAttemptRepository) FetchOneByUnitIDAndUserID(ctx context.Context, userID string, unit string) (exam_management.ExamManagement, error) {
	//TODO implement me
	panic("implement me")
}

//func (u *userAttemptRepository) FetchOneByUnitIDAndUserID(ctx context.Context, userID string, unit string) (user_attempt_domain.UserProcess, error) {
//	//collectionExam := u.database.Collection(u.collectionExam)
//	//collectionExercise := u.database.Collection(u.collectionExercise)
//	//collectionUserAttempt := u.database.Collection(u.collectionUserAttempt)
//	//collectionQuiz := u.database.Collection(u.collectionQuiz)
//	//
//	//idUser, err := primitive.ObjectIDFromHex(userID)
//	//if err != nil {
//	//	return user_attempt_domain.UserProcess{}, err
//	//}
//	//
//	//idUnit, err := primitive.ObjectIDFromHex(unit)
//	//if err != nil {
//	//	return user_attempt_domain.UserProcess{}, err
//	//}
//	//
//	//filter := bson.M{"user_id": idUser, "unit_id": idUnit}
//	//panic(error())
//
//}

func (u *userAttemptRepository) UpdateExamManagementByExamID(ctx context.Context, userID exam_management.ExamManagement) error {
	collectionUserAttempt := u.database.Collection(u.collectionUserAttempt)

	filter := bson.D{{Key: "user_id", Value: userID.UserID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"exam_id":        userID.ExamID,
		"score":          userID.Score,
		"process_status": userID.ProcessStatus,
		"completed_date": userID.CompletedDate,
		"updated_at":     userID.UpdatedAt,
	}}}

	_, err := collectionUserAttempt.UpdateOne(ctx, filter, &update)
	return err
}

func (u *userAttemptRepository) UpdateExamManagementByQuizID(ctx context.Context, userID exam_management.ExamManagement) error {
	collectionUserAttempt := u.database.Collection(u.collectionUserAttempt)

	filter := bson.D{{Key: "user_id", Value: userID.UserID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"quiz_id":        userID.ExamID,
		"score":          userID.Score,
		"process_status": userID.ProcessStatus,
		"completed_date": userID.CompletedDate,
		"updated_at":     userID.UpdatedAt,
	}}}

	_, err := collectionUserAttempt.UpdateOne(ctx, filter, &update)
	return err
}

func (u *userAttemptRepository) FetchManyByUserID(ctx context.Context, userID string) (exam_management.Response, error) {
	collectionUserAttempt := u.database.Collection(u.collectionUserAttempt)

	idUser, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return exam_management.Response{}, err
	}

	filter := bson.M{"user_id": idUser}

	var examManagement exam_management.ExamManagement
	err = collectionUserAttempt.FindOne(ctx, filter).Decode(examManagement)
	if err != nil {
		return exam_management.Response{}, err
	}

	userAttemptRes := exam_management.Response{
		ExamManagement: examManagement,
	}
	return userAttemptRes, nil
}

func (u *userAttemptRepository) CreateExamManagementByExerciseID(ctx context.Context, userID exam_management.ExamManagement) error {
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
	if err := checkExistence(collectionExam, userID.ExamID, "exam"); err != nil {
		return err
	}

	// Kiểm tra sự tồn tại của Quiz
	if err := checkExistence(collectionQuiz, userID.QuizID, "quiz"); err != nil {
		return err
	}

	// Kiểm tra sự tồn tại của Exercise
	if err := checkExistence(collectionExercise, userID.ExerciseID, "exercise"); err != nil {
		return err
	}

	// Chèn userProcess vào collectionUserAttempt
	_, err := collectionUserAttempt.InsertOne(ctx, userID)
	return err
}

func (u *userAttemptRepository) UpdateExamManagementByUserID(ctx context.Context, userID exam_management.ExamManagement) error {
	collectionUserAttempt := u.database.Collection(u.collectionUserAttempt)

	filter := bson.D{{Key: "user_id", Value: userID.UserID}}
	update := bson.D{{Key: "$set", Value: bson.M{
		"score":          userID.Score,
		"process_status": userID.ProcessStatus,
		"completed_date": userID.CompletedDate,
		"updated_at":     userID.UpdatedAt,
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
