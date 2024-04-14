package user_attempt

import (
	"clean-architecture/domain/user_attempt"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userAttemptRepository struct {
	database              mongo.Database
	collectionLesson      string
	collectionUnit        string
	collectionVocabulary  string
	collectionUserAttempt string
}

func (u *userAttemptRepository) FetchManyByUserID(c context.Context) (user_attempt_domain.Response, error) {
	collectionExam := u.database.Collection(u.collectionUserAttempt)

	cursor, err := collectionExam.Find(c, bson.D{})
	if err != nil {
		return user_attempt_domain.Response{}, err
	}

	var userAttempts []user_attempt_domain.UserProcess
	for cursor.Next(c) {
		var userAttempt user_attempt_domain.UserProcess
		if err = cursor.Decode(&userAttempt); err != nil {
			return user_attempt_domain.Response{}, err
		}

		// Thêm user_attempt vào slice lessons
		userAttempts = append(userAttempts, userAttempt)
	}

	userAttemptRes := user_attempt_domain.Response{
		UserProcess: userAttempts,
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

func (u *userAttemptRepository) DeleteAllByUserID(c context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}

func NewUserAttemptRepository(db mongo.Database, collectionUserAttempt string, collectionLesson string, collectionUnit string, collectionVocabulary string) user_attempt_domain.IUserProcessRepository {
	return &userAttemptRepository{
		database:              db,
		collectionUserAttempt: collectionUserAttempt,
		collectionLesson:      collectionLesson,
		collectionUnit:        collectionUnit,
		collectionVocabulary:  collectionVocabulary,
	}
}
