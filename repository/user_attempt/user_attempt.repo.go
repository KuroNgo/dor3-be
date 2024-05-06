package user_attempt

import (
	"clean-architecture/domain/user_attempt"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userAttemptRepository struct {
	database              mongo.Database
	collectionExam        string
	collectionExercise    string
	collectionQuiz        string
	collectionUserAttempt string
}

func NewUserAttemptRepository(db mongo.Database, collectionUserAttempt string, collectionLesson string, collectionUnit string, collectionVocabulary string) user_attempt_domain.IUserProcessRepository {
	return &userAttemptRepository{
		database:              db,
		collectionUserAttempt: collectionUserAttempt,
		collectionExam:        collectionLesson,
		collectionExercise:    collectionUnit,
		collectionQuiz:        collectionVocabulary,
	}
}

func (u *userAttemptRepository) FetchManyByUserID(c context.Context) (user_attempt_domain.Response, error) {
	collectionExam := u.database.Collection(u.collectionUserAttempt)

	var userAttempt user_attempt_domain.UserProcess
	err := collectionExam.FindOne(c, bson.D{}).Decode(userAttempt)
	if err != nil {
		return user_attempt_domain.Response{}, err
	}

	//statistics := user_attempt_domain.Statistics{
	//
	//}
	//
	userAttemptRes := user_attempt_domain.Response{
		UserProcess: userAttempt,
	}
	return userAttemptRes, nil

}

func (u *userAttemptRepository) CreateOneByUserID(c context.Context, userID string) error {
	//collectionUserAttempt := u.database.Collection(u.collectionUserAttempt)
	//collectionLesson := u.database.Collection(u.collectionLesson)
	//collectionUnit := u.database.Collection(u.collectionUnit)
	//collectionVocabulary := u.database.Collection(u.collectionVocabulary)
	//
	//filter := bson.M{"user_id": userID}
	//TODO implement me
	panic("implement me")
}

func (u *userAttemptRepository) UpdateAttemptByUserID(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (u *userAttemptRepository) DeleteAllByUserID(c context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}
