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

func (i *imageRepository) CreateMany(ctx context.Context, images []*image_domain.Image) error {
	collection := i.database.Collection(i.collection)

	var imageNames []string
	for _, img := range images {
		imageNames = append(imageNames, img.ImageName)
	}

	filter := bson.M{"image-name": bson.M{"$in": imageNames}}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		errors.New("one or more image names already exist")
	}

	var documents []interface{}
	for _, img := range images {
		documents = append(documents, img)
	}

	_, err = collection.InsertMany(ctx, documents)
	if err != nil {
		return err
	}

	return nil
}

func (i *imageRepository) GetURLByName(ctx context.Context, name string) (image_domain.Image, error) {
	collection := i.database.Collection(i.collection)
	var image image_domain.Image
	err := collection.FindOne(ctx, bson.M{"image_name": name}).Decode(&image)
	return image, err
}

func (i *imageRepository) FetchMany(ctx context.Context) ([]image_domain.Image, error) {
	collection := i.database.Collection(i.collection)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var image []image_domain.Image

	err = cursor.All(ctx, &image)
	if image == nil {
		return []image_domain.Image{}, err
	}

	return image, err
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

func NewImageRepository(db mongo.Database, collection string) image_domain.IImageRepository {
	return &imageRepository{
		database:   db,
		collection: collection,
	}
}
