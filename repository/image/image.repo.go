package image_repository

import (
	image_domain "clean-architecture/domain/image"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"sync"
)

type imageRepository struct {
	database   *mongo.Database
	collection string
}

func (i *imageRepository) FetchByCategory(ctx context.Context, category string, page string) (image_domain.Response, error) {
	collectionImage := i.database.Collection(i.collection)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return image_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	calculate := make(chan int64)

	go func() {
		defer close(calculate)
		// Đếm tổng số lượng tài liệu trong collection
		count, err := collectionImage.CountDocuments(ctx, bson.D{})
		if err != nil {
			return
		}

		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)
		if cal2 != 0 {
			calculate <- cal1
		}
	}()

	filter := bson.M{"category": category}
	cursor, err := collectionImage.Find(ctx, filter, findOptions)
	if err != nil {
		return image_domain.Response{}, err
	}

	var images []image_domain.Image
	for cursor.Next(ctx) {
		var image image_domain.Image
		if err := cursor.Decode(&image); err != nil {
			return image_domain.Response{}, err
		}

		image.Category = category
		images = append(images, image)
	}

	cal := <-calculate

	imgRes := image_domain.Response{
		Page:  cal,
		Image: images,
	}
	return imgRes, nil
}

func NewImageRepository(db *mongo.Database, collection string) image_domain.IImageRepository {
	return &imageRepository{
		database:   db,
		collection: collection,
	}
}

func (i *imageRepository) CreateMany(ctx context.Context, images []*image_domain.Image) error {
	collection := i.database.Collection(i.collection)

	var documents []interface{}
	for _, image := range images {
		filter := bson.M{"image_name": image.ImageName, "image_url": image.ImageUrl}

		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("the image with name '%s' and URL '%s' already exists", image.ImageName, image.ImageUrl)
		}

		documents = append(documents, image)
	}

	_, err := collection.InsertMany(ctx, documents)
	return err
}

func (i *imageRepository) GetURLByName(ctx context.Context, name string) (image_domain.Image, error) {
	collection := i.database.Collection(i.collection)
	var image image_domain.Image
	err := collection.FindOne(ctx, bson.M{"image_name": name}).Decode(&image)
	return image, err
}

func (i *imageRepository) FetchMany(ctx context.Context, page string) (image_domain.Response, error) {
	collection := i.database.Collection(i.collection)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return image_domain.Response{}, errors.New("invalid page number")
	}

	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	var count, size int64
	var images []image_domain.Image
	var mu sync.Mutex

	// Đồng bộ hóa việc đếm số lượng tài liệu và tính toán
	mu.Lock()
	defer mu.Unlock()

	// Đếm tổng số lượng tài liệu trong collection
	count, err = collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return image_domain.Response{}, fmt.Errorf("error counting documents: %v", err)
	}

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return image_domain.Response{}, fmt.Errorf("error finding documents: %v", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	go func() {
		for cursor.Next(ctx) {
			var doc image_domain.Image
			if err := cursor.Decode(&doc); err != nil {
				log.Println("Error decoding document:", err)
				continue
			}
			size += doc.Size
			images = append(images, doc)
		}
	}()

	statistics := image_domain.Statistics{
		Count:           count,
		MaxSizeMB:       1024 / 1024,
		MaxSizeKB:       1024 * 1024,
		SizeKB:          size,
		SizeMB:          size / 1024,
		SizeRemainingKB: (1024 * 1024) - size,
		SizeRemainingMB: (1024 * 1024 * 1024) - size,
	}

	response := image_domain.Response{
		Image:      images,
		Statistics: statistics,
		Page:       count/int64(perPage) + 1,
	}

	return response, nil
}

func (i *imageRepository) UpdateOne(ctx context.Context, imageID string, updatedImage *image_domain.Image) error {
	collection := i.database.Collection(i.collection)

	filter := bson.M{"_id": imageID}
	update := bson.M{"$set": updatedImage}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update image: %w", err)
	}

	return nil
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
