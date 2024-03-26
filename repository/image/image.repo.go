package image_repository

import (
	image_domain "clean-architecture/domain/image"
	"clean-architecture/infrastructor/mongo"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type imageRepository struct {
	database   mongo.Database
	collection string
}

func NewImageRepository(db mongo.Database, collection string) image_domain.IImageRepository {
	return &imageRepository{
		database:   db,
		collection: collection,
	}
}
func (i *imageRepository) GetURLByName(ctx context.Context, name string) (image_domain.Image, error) {
	collection := i.database.Collection(i.collection)
	var image image_domain.Image
	err := collection.FindOne(ctx, bson.M{"image_name": name}).Decode(&image)
	return image, err
}

func (i *imageRepository) FetchMany(ctx context.Context) (image_domain.Response, error) {
	collection := i.database.Collection(i.collection)

	// Đếm tổng số lượng tài liệu trong collection
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return image_domain.Response{}, err
	}

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return image_domain.Response{}, err
	}

	// Lặp qua các tài liệu và tính tổng của trường FieldToSum
	var size int64
	var images []image_domain.Image
	for cursor.Next(context.Background()) {
		var doc image_domain.Image
		if err = cursor.Decode(&doc); err != nil {
			return image_domain.Response{}, errors.New("")
		}
		images = append(images, doc)
		size += doc.Size
	}

	var image []image_domain.Image

	err = cursor.All(ctx, &image)
	if image == nil {
		return image_domain.Response{}, err
	}
	// Tạo cấu trúc dữ liệu Response
	response := image_domain.Response{
		Image: image,
		Count: count,
		Size:  size,
	}

	return response, err
}

func (i *imageRepository) UpdateOne(ctx context.Context, imageID string, image image_domain.Image) error {
	collection := i.database.Collection(i.collection)
	objID, err := primitive.ObjectIDFromHex(imageID)

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{}}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func (i *imageRepository) CreateOne(ctx context.Context, image *image_domain.Image) error {
	collection := i.database.Collection(i.collection)

	filter := bson.M{"image_name": image.ImageName, "image_url": image.ImageUrl}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the image name did exists")
	}

	_, err = collection.InsertOne(ctx, image)
	return err
}

func (i *imageRepository) DeleteOne(ctx context.Context, imageID string) error {
	collection := i.database.Collection(i.collection)
	objID, err := primitive.ObjectIDFromHex(imageID)
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

func (i *imageRepository) DeleteMany(ctx context.Context, imageID ...string) error {
	collection := i.database.Collection(i.collection)
	var objIDs []primitive.ObjectID

	for _, audioID := range imageID {
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
		return errors.New("the images do not exists or has been removed")
	}
	_, err = collection.DeleteMany(ctx, filter)
	return err
}
