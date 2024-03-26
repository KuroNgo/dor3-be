package audio_repository

import (
	audio_domain "clean-architecture/domain/audio"
	"clean-architecture/infrastructor/mongo"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type audioRepository struct {
	database   mongo.Database
	collection string
}

func NewAudioRepository(db mongo.Database, collection string) audio_domain.IAudioRepository {
	return &audioRepository{
		database:   db,
		collection: collection,
	}
}

func (a *audioRepository) FetchMany(ctx context.Context) (audio_domain.Response, error) {
	collection := a.database.Collection(a.collection)

	// Đếm tổng số lượng tài liệu trong collection
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return audio_domain.Response{}, err
	}

	// Sử dụng Find() để lấy dữ liệu từ collection
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return audio_domain.Response{}, err
	}
	defer func(cursor mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var audio []audio_domain.Audio

	// Decode dữ liệu từ cursor vào slice audio
	if err = cursor.All(ctx, &audio); err != nil {
		return audio_domain.Response{}, err
	}

	// Tạo cấu trúc dữ liệu Response
	response := audio_domain.Response{
		Audio: audio,
		Count: count,
	}

	return response, nil
}

func (a *audioRepository) CreateOne(ctx context.Context, audio *audio_domain.Audio) error {
	collection := a.database.Collection(a.collection)

	filter := bson.M{"filename": audio.Filename}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the audio id did exists")
	}

	_, err = collection.InsertOne(ctx, audio)
	return err
}

func (a *audioRepository) UpdateOne(ctx context.Context, audioID string, audio audio_domain.Audio) error {
	collection := a.database.Collection(a.collection)
	objID, err := primitive.ObjectIDFromHex(audioID)

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (a *audioRepository) DeleteOne(ctx context.Context, audioID string) error {
	collection := a.database.Collection(a.collection)
	objID, err := primitive.ObjectIDFromHex(audioID)
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
	if count <= 0 {
		return errors.New(`the course had been removed or does not exists`)
	}
	_, err = collection.DeleteOne(ctx, filter)
	return err
}

func (a *audioRepository) DeleteMany(ctx context.Context, audioIDs ...string) error {
	collection := a.database.Collection(a.collection)
	var objIDs []primitive.ObjectID

	for _, audioID := range audioIDs {
		objID, err := primitive.ObjectIDFromHex(audioID)
		if err != nil {
			return err
		}
		objIDs = append(objIDs, objID)
	}

	filter := bson.M{
		"_id": bson.M{"$in": objIDs}, // use $in operator for delete many document in the same time
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count <= 0 {
		return errors.New("the audios do not exists or had been removed")
	}
	_, err = collection.DeleteMany(ctx, filter)
	return err
}
