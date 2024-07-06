package image_repository

import (
	image_domain "clean-architecture/domain/image"
	"clean-architecture/internal"
	"clean-architecture/internal/cache/memory"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"sync"
	"time"
)

type imageRepository struct {
	database   *mongo.Database
	collection string
}

func NewImageRepository(db *mongo.Database, collection string) image_domain.IImageRepository {
	return &imageRepository{
		database:   db,
		collection: collection,
	}
}

var (
	imageResCache = memory.NewTTL[string, image_domain.Response]()
	imageCache    = memory.NewTTL[string, image_domain.Image]()

	wg           sync.WaitGroup
	mu           sync.Mutex
	isProcessing bool
)

const (
	cacheTTL = 5 * time.Minute
)

func (i *imageRepository) FetchByCategory(ctx context.Context, category string, page string) (image_domain.Response, error) {
	errCh := make(chan error, 1)
	imageResCh := make(chan image_domain.Response, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := imageResCache.Get(category + page)
		if found {
			imageResCh <- data
		}
	}()

	go func() {
		defer close(imageResCh)
		wg.Wait()
	}()

	imageResData := <-imageResCh
	if !internal.IsZeroValue(imageResData) {
		return imageResData, nil
	}

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
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var images []image_domain.Image
	images = make([]image_domain.Image, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var image image_domain.Image
		if err = cursor.Decode(&image); err != nil {
			return image_domain.Response{}, err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			image.Category = category
			images = append(images, image)
		}()
	}
	wg.Wait()

	cal := <-calculate
	statistics, _ := i.Statistics(ctx)

	imgRes := image_domain.Response{
		Page:        cal,
		CurrentPage: pageNumber,
		Statistics:  statistics,
		Image:       images,
	}

	imageResCache.Set(category+page, imgRes, cacheTTL)

	select {
	case err = <-errCh:
		return image_domain.Response{}, err
	default:
		return imgRes, nil
	}
}

func (i *imageRepository) GetURLByName(ctx context.Context, name string) (image_domain.Image, error) {
	imageCh := make(chan image_domain.Image, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := imageCache.Get(name)
		if found {
			imageCh <- data
		}
	}()

	go func() {
		defer close(imageCh)
		wg.Wait()
	}()

	imageData := <-imageCh
	if !internal.IsZeroValue(imageData) {
		return imageData, nil
	}

	collection := i.database.Collection(i.collection)
	var image image_domain.Image
	err := collection.FindOne(ctx, bson.M{"image_name": name}).Decode(&image)

	imageCache.Set(name, image, cacheTTL)
	return image, err
}

func (i *imageRepository) FetchMany(ctx context.Context, page string) (image_domain.Response, error) {
	errCh := make(chan error, 1)
	imageResCh := make(chan image_domain.Response, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, found := imageResCache.Get(page)
		if found {
			imageResCh <- data
		}
	}()

	go func() {
		defer close(imageResCh)
		wg.Wait()
	}()

	imageResData := <-imageResCh
	if !internal.IsZeroValue(imageResData) {
		return imageResData, nil
	}

	collection := i.database.Collection(i.collection)

	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		return image_domain.Response{}, errors.New("invalid page number")
	}

	const perPage = 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	var count int64

	// Count the total number of documents in the collection
	count, err = collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return image_domain.Response{}, fmt.Errorf("error counting documents: %v", err)
	}

	// Fetch the current page of images
	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return image_domain.Response{}, fmt.Errorf("error finding documents: %v", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var images []image_domain.Image
	images = make([]image_domain.Image, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var doc image_domain.Image
		if err := cursor.Decode(&doc); err != nil {
			return image_domain.Response{}, fmt.Errorf("error decoding document: %v", err)
		}

		wg.Add(1)
		go func(doc image_domain.Image) {
			defer wg.Done()
			images = append(images, doc)
		}(doc)
	}
	wg.Wait()

	if err = cursor.Err(); err != nil {
		return image_domain.Response{}, fmt.Errorf("error iterating cursor: %v", err)
	}

	totalPages := (count + int64(perPage) - 1) / int64(perPage)
	statistics, _ := i.Statistics(ctx)

	response := image_domain.Response{
		Image:      images,
		Statistics: statistics,
		Page:       totalPages,
	}

	imageResCache.Set(page, response, cacheTTL)

	select {
	case err = <-errCh:
		return image_domain.Response{}, err
	default:
		return response, nil
	}

}

func (i *imageRepository) CreateMany(ctx context.Context, images []*image_domain.Image) error {
	// Khóa lock giúp bảo vệ course
	mu.Lock()
	defer mu.Unlock()

	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	isProcessing = true
	defer func() {
		isProcessing = false
	}()

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

func (i *imageRepository) UpdateOne(ctx context.Context, updatedImage *image_domain.Image) error {
	// Khóa lock giúp bảo vệ course
	mu.Lock()
	defer mu.Unlock()

	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	isProcessing = true
	defer func() {
		isProcessing = false
	}()

	collection := i.database.Collection(i.collection)

	idImage, _ := primitive.ObjectIDFromHex(updatedImage.Id.Hex())
	filter := bson.M{"_id": idImage}
	update := bson.M{"$set": bson.M{
		"image_name": updatedImage.ImageName,
		"image_url":  updatedImage.ImageUrl,
		"asset_id":   updatedImage.AssetId,
		"size":       updatedImage.Size,
	}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update image: %w", err)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		imageCache.Clear()
	}()

	go func() {
		defer wg.Done()
		imageResCache.Clear()
	}()
	wg.Wait()
	return nil
}

func (i *imageRepository) CreateOne(ctx context.Context, image *image_domain.Image) error {
	// Khóa lock giúp bảo vệ course
	mu.Lock()
	defer mu.Unlock()

	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	isProcessing = true
	defer func() {
		isProcessing = false
	}()

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
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		imageResCache.Clear()
	}()
	wg.Wait()

	return nil
}

func (i *imageRepository) DeleteOne(ctx context.Context, imageID string) error {
	// Khóa lock giúp bảo vệ course
	mu.Lock()
	defer mu.Unlock()

	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	isProcessing = true
	defer func() {
		isProcessing = false
	}()

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
	if err != nil {
		return err
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		imageCache.Remove(imageID)
	}()

	go func() {
		defer wg.Done()
		imageResCache.Clear()
	}()
	wg.Wait()

	return nil
}

func (i *imageRepository) DeleteMany(ctx context.Context, imageID ...string) error {
	// Khóa lock giúp bảo vệ course
	mu.Lock()
	defer mu.Unlock()

	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	isProcessing = true
	defer func() {
		isProcessing = false
	}()

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
	if err != nil {
		return err
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		imageCache.Clear()
	}()

	go func() {
		defer wg.Done()
		imageResCache.Clear()
	}()
	wg.Wait()

	return nil
}

func (i *imageRepository) Statistics(ctx context.Context) (image_domain.Statistics, error) {
	collectionImage := i.database.Collection(i.collection)

	const (
		MaxSizeMB   = 1024.0
		MaxSizeKB   = 1024 * 1024
		TotalSizeKB = 1024 * 1024 // 1GB in KB
		TotalSizeMB = 1024        // 1GB in MB
	)

	count, err := collectionImage.CountDocuments(ctx, bson.D{})
	if err != nil {
		return image_domain.Statistics{}, err
	}

	// Use an aggregation pipeline to calculate the total size
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{{"_id", nil}, {"totalSize", bson.D{{"$sum", "$size"}}}}}},
	}

	var result struct {
		TotalSize int64 `bson:"totalSize"`
	}

	cursor, err := collectionImage.Aggregate(ctx, pipeline)
	if err != nil {
		return image_domain.Statistics{}, fmt.Errorf("error aggregating documents: %v", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return image_domain.Statistics{}, fmt.Errorf("error decoding aggregation result: %v", err)
		}
	}

	if err := cursor.Err(); err != nil {
		return image_domain.Statistics{}, fmt.Errorf("error iterating cursor: %v", err)
	}

	totalSize := result.TotalSize

	statistics := image_domain.Statistics{
		Count:           count,
		MaxSizeMB:       MaxSizeMB,
		MaxSizeKB:       MaxSizeKB,
		SizeKB:          totalSize,
		SizeMB:          totalSize / 1024,
		SizeRemainingKB: TotalSizeKB - totalSize,
		SizeRemainingMB: (TotalSizeKB - totalSize) / 1024,
	}

	return statistics, nil
}
