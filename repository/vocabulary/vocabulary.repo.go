package vocabulary_repository

import (
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/internal"
	"clean-architecture/internal/cache"
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// vocabularyRepository defines the structure for the repository that interacts with the vocabulary collections in the database.
type vocabularyRepository struct {
	database             *mongo.Database // Reference to the MongoDB database
	collectionVocabulary string          // Name of the vocabulary collection
	collectionMark       string          // Name of the mark collection
	collectionUnit       string          // Name of the unit collection
	collectionLesson     string          // Name of the lesson collection
}

// NewVocabularyRepository creates a new instance of vocabularyRepository.
// It initializes the repository with the provided MongoDB database and collection names.
// Parameters:
//   - db: The MongoDB database.
//   - collectionVocabulary: The name of the vocabulary collection.
//   - collectionMark: The name of the mark collection.
//   - collectionUnit: The name of the unit collection.
//   - collectionLesson: The name of the lesson collection.
//
// Returns:
//   - vocabulary_domain.IVocabularyRepository: An interface to interact with the vocabulary repository.
func NewVocabularyRepository(db *mongo.Database, collectionVocabulary string, collectionMark string, collectionUnit string, collectionLesson string) vocabulary_domain.IVocabularyRepository {
	return &vocabularyRepository{
		database:             db,
		collectionVocabulary: collectionVocabulary,
		collectionMark:       collectionMark,
		collectionUnit:       collectionUnit,
		collectionLesson:     collectionLesson,
	}
}

// Cache variables and synchronization primitives
var (
	vocabularyCache         = cache.NewTTL[string, vocabulary_domain.Vocabulary]()        // Cache for single vocabulary entries
	vocabulariesCache       = cache.NewTTL[string, []vocabulary_domain.Vocabulary]()      // Cache for multiple vocabulary entries
	vocabulariesSearchCache = cache.NewTTL[string, vocabulary_domain.SearchingResponse]() // Cache for vocabulary search responses
	vocabularyResponseCache = cache.NewTTL[string, vocabulary_domain.Response]()          // Cache for vocabulary response pages
	vocabularyPrimOIDCache  = cache.NewTTL[string, primitive.ObjectID]()                  // Cache for vocabulary ObjectID mappings
	vocabularyArrCache      = cache.NewTTL[string, []string]()                            // Cache for arrays of vocabulary strings

	mu           sync.Mutex     // Mutex for ensuring thread safety
	wg           sync.WaitGroup // WaitGroup for managing goroutines
	isProcessing bool           // Flag to indicate if a process is ongoing
)

// FindVocabularyIDByVocabularyConfigInAdmin finds the ObjectID of a vocabulary entry by its configuration word within an admin context.
// It first attempts to retrieve the ID from the cache, and if not found, queries the database and caches the result.
// Parameters:
//   - ctx: The context for managing deadlines, cancellation signals, and other request-scoped values.
//   - word: The configuration word of the vocabulary entry to be searched.
//
// Returns:
//   - primitive.ObjectID: The ObjectID of the vocabulary entry.
//   - error: An error message if the fetch or conditions are not met.
func (v *vocabularyRepository) FindVocabularyIDByVocabularyConfigInAdmin(ctx context.Context, word string) (primitive.ObjectID, error) {
	vocabularyCh := make(chan primitive.ObjectID, 1) // Channel to handle the vocabulary ObjectID

	wg.Add(1) // Add a wait group to manage goroutines
	go func() {
		defer wg.Done()
		data, found := vocabularyPrimOIDCache.Get(word)
		if found {
			vocabularyCh <- data
		}
	}()

	go func() {
		defer close(vocabularyCh) // Ensure vocabularyCh is closed after processing
		wg.Wait()                 // Wait for all goroutines to finish
	}()

	vocabularyData := <-vocabularyCh // Receive data from the channel
	if !internal.IsZeroValue(vocabularyData) {
		return vocabularyData, nil
	}

	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	filter := bson.M{"word_for_config": word} // Filter to find the vocabulary entry by the configuration word
	var data struct {
		Id primitive.ObjectID `bson:"_id"` // Structure to hold the ObjectID of the found document
	}

	// Query the database for the vocabulary entry
	err := collectionVocabulary.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// Cache the found ObjectID
	vocabularyPrimOIDCache.Set(word, data.Id, 5*time.Minute)
	return data.Id, nil
}

// GetLatestVocabularyInAdmin retrieves the latest vocabulary entries created within the last 24 hours in an admin context.
// It first attempts to retrieve the data from the cache, and if not found, queries the database and caches the result.
// Parameters:
//   - ctx: The context for managing deadlines, cancellation signals, and other request-scoped values.
//
// Returns:
//   - []string: A list of the latest vocabulary words.
//   - error: An error message if the fetch or conditions are not met.
func (v *vocabularyRepository) GetLatestVocabularyInAdmin(ctx context.Context) ([]string, error) {
	errCh := make(chan error, 1)              // Channel to handle errors
	vocabularyArrCh := make(chan []string, 1) // Channel to handle the vocabulary array response

	wg.Add(1) // Add a wait group to manage goroutines
	go func() {
		defer wg.Done()
		data, found := vocabularyArrCache.Get("latest")
		if found {
			vocabularyArrCh <- data
		}
	}()

	go func() {
		defer close(vocabularyArrCh) // Ensure vocabularyArrCh is closed after processing
		wg.Wait()                    // Wait for all goroutines to finish
	}()

	vocabularyArrData := <-vocabularyArrCh // Receive data from the channel
	if !internal.IsZeroValue(vocabularyArrData) {
		return vocabularyArrData, nil
	}

	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	var vocabularies []string
	yesterday := time.Now().Add(-24 * time.Hour)             // Get the timestamp for 24 hours ago
	filter := bson.M{"created_at": bson.M{"$gt": yesterday}} // Filter for documents created in the last 24 hours

	// Query the database for the latest vocabulary entries
	cursor, err := collectionVocabulary.Find(ctx, filter, options.Find().SetSort(bson.D{{"_id", -1}}))
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	// Iterate over the cursor to process the results
	for cursor.Next(ctx) {
		var result bson.M
		if err = cursor.Decode(&result); err != nil {
			return nil, err
		}

		wg.Add(1)
		go func(result bson.M) {
			defer wg.Done()
			word, ok := result["word"].(string)
			if !ok {
				errCh <- errors.New("failed to parse word from result")
				return
			}
			vocabularies = append(vocabularies, word)
		}(result)
	}
	wg.Wait()

	// Cache the retrieved vocabulary list
	vocabularyArrCache.Set("latest", vocabularies, 5*time.Minute)

	// Check for any errors that might have occurred during the goroutines
	select {
	case err = <-errCh:
		return nil, err
	default:
		return vocabularies, nil
	}
}

// GetVocabularyByIdInAdmin retrieves a vocabulary entry by its ID within an admin context.
// It first attempts to retrieve the data from the cache, and if not found, queries the database and caches the result.
// Parameters:
//   - ctx: The context for managing deadlines, cancellation signals, and other request-scoped values.
//   - id: The ID of the vocabulary entry to be retrieved.
//
// Returns:
//   - vocabulary_domain.Vocabulary: The vocabulary entry corresponding to the given ID.
//   - error: An error message if the fetch or conditions are not met.
func (v *vocabularyRepository) GetVocabularyByIdInAdmin(ctx context.Context, id string) (vocabulary_domain.Vocabulary, error) {
	vocabularyCh := make(chan vocabulary_domain.Vocabulary, 1) // Channel to handle the vocabulary response

	wg.Add(1) // Add a wait group to manage goroutines
	go func() {
		defer wg.Done()
		data, found := vocabularyCache.Get(id)
		if found {
			vocabularyCh <- data
		}
	}()

	go func() {
		defer close(vocabularyCh) // Ensure vocabularyCh is closed after processing
		wg.Wait()                 // Wait for all goroutines to finish
	}()

	vocabularyData := <-vocabularyCh // Receive data from the channel
	if !internal.IsZeroValue(vocabularyData) {
		return vocabularyData, nil
	}

	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	// Convert the string ID to a MongoDB ObjectID
	idVocabulary, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return vocabulary_domain.Vocabulary{}, err
	}

	filter := bson.M{"_id": idVocabulary} // Filter to find the vocabulary entry by ID

	var vocabulary vocabulary_domain.Vocabulary
	// Query the database for the vocabulary entry
	err = collectionVocabulary.FindOne(ctx, filter).Decode(&vocabulary)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return vocabulary_domain.Vocabulary{}, err
	}

	// Cache the retrieved vocabulary entry
	vocabularyCache.Set(id, vocabulary, 5*time.Minute)
	return vocabulary, nil
}

// GetAllVocabularyInAdmin retrieves all vocabulary words from the database in an admin context.
// It first attempts to retrieve the data from the cache, and if not found, queries the database and caches the result.
// Parameters:
//   - ctx: The context for managing deadlines, cancellation signals, and other request-scoped values.
//
// Returns:
//   - []string: A list of all vocabulary words.
//   - error: An error message if the fetch or conditions are not met.
func (v *vocabularyRepository) GetAllVocabularyInAdmin(ctx context.Context) ([]string, error) {
	errCh := make(chan error)           // Channel to handle errors
	vocabularyCh := make(chan []string) // Channel to handle the vocabulary array response

	wg.Add(1) // Add a wait group to manage goroutines
	go func() {
		defer wg.Done()
		data, found := vocabularyArrCache.Get("all")
		if found {
			vocabularyCh <- data
		}
	}()

	go func() {
		defer close(vocabularyCh) // Ensure vocabularyCh is closed after processing
		wg.Wait()                 // Wait for all goroutines to finish
	}()

	vocabularyData := <-vocabularyCh // Receive data from the channel
	if !internal.IsZeroValue(vocabularyData) {
		return vocabularyData, nil
	}

	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	var vocabularies []string

	// Query the database to fetch all documents from the vocabulary collection
	cursor, err := collectionVocabulary.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	// Iterate over the cursor to process each document
	for cursor.Next(ctx) {
		var result bson.M
		if err = cursor.Decode(&result); err != nil {
			return nil, err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			word, ok := result["word"].(string)
			if !ok {
				errCh <- errors.New("failed to parse word from result")
				return
			}
			vocabularies = append(vocabularies, word)
		}()
	}
	wg.Wait()

	// Cache the retrieved vocabulary list for future queries
	vocabularyArrCache.Set("all", vocabularies, 5*time.Minute)

	// Check for any errors that might have occurred during the goroutines
	select {
	case err = <-errCh:
		return nil, err
	default:
		return vocabularies, nil
	}
}

// FetchByIdUnitInAdmin fetches vocabulary entries by unit ID in an admin context.
// It first attempts to retrieve the data from the cache, and if not found, queries the database and caches the result.
// Parameters:
//   - ctx: The context for managing deadlines, cancellation signals, and other request-scoped values.
//   - idUnit: The ID of the unit whose vocabulary entries are to be fetched.
//
// Returns:
//   - []vocabulary_domain.Vocabulary: A list of vocabulary entries belonging to the specified unit ID.
//   - error: An error message if the fetch or conditions are not met.
func (v *vocabularyRepository) FetchByIdUnitInAdmin(ctx context.Context, idUnit string) ([]vocabulary_domain.Vocabulary, error) {
	errCh := make(chan error, 1)                                   // Channel to handle errors
	vocabulariesCh := make(chan []vocabulary_domain.Vocabulary, 1) // Channel to handle vocabulary array response

	wg.Add(1) // Add a wait group to manage goroutines
	go func() {
		defer wg.Done()
		data, found := vocabulariesCache.Get(idUnit)
		if found {
			vocabulariesCh <- data
		}
	}()

	go func() {
		defer close(vocabulariesCh) // Ensure vocabulariesCh is closed after processing
		wg.Wait()                   // Wait for all goroutines to finish
	}()

	vocabulariesData := <-vocabulariesCh // Receive data from the channel
	if !internal.IsZeroValue(vocabulariesData) {
		return vocabulariesData, nil
	}

	// Get references to the vocabulary and unit collections from the database
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	// Convert the string ID to a MongoDB ObjectID for the unit
	unitID, err := primitive.ObjectIDFromHex(idUnit)
	if err != nil {
		return nil, fmt.Errorf("invalid unit id: %w", err)
	}

	// Query the unit collection to fetch the unit corresponding to the given ID
	filterUnit := bson.M{"_id": unitID}
	var unit unit_domain.Unit
	err = collectionUnit.FindOne(ctx, filterUnit).Decode(&unit)
	if err != nil {
		return nil, fmt.Errorf("failed to find unit: %w", err)
	}

	// Filter to find vocabulary entries belonging to the unit
	filter := bson.M{"unit_id": unit.ID}
	cursor, err := collectionVocabulary.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find vocabularies: %w", err)
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			errCh <- err
			return
		}
	}()

	// Slice to hold the vocabulary results
	var vocabularies []vocabulary_domain.Vocabulary

	// Iterate over the cursor to process each vocabulary entry
	for cursor.Next(ctx) {
		var vocabulary vocabulary_domain.Vocabulary
		if err = cursor.Decode(&vocabulary); err != nil {
			return nil, fmt.Errorf("failed to decode vocabulary: %w", err)
		}

		wg.Add(1)
		go func(vocabulary vocabulary_domain.Vocabulary) {
			defer wg.Done()
			// Append the vocabulary entry to the slice
			vocabularies = append(vocabularies, vocabulary)
		}(vocabulary)
	}
	wg.Wait()

	// Cache the retrieved vocabulary entries for future queries
	vocabulariesCache.Set(idUnit, vocabularies, 5*time.Minute)

	// Check for any errors that might have occurred during the goroutines
	select {
	case err = <-errCh:
		return nil, err
	default:
		return vocabularies, nil
	}
}

// FetchByWordInBoth searches for vocabulary entries by word pattern in an admin context.
// It first attempts to retrieve the data from the cache, and if not found, queries the database and caches the result.
// Parameters:
//   - ctx: The context for managing deadlines, cancellation signals, and other request-scoped values.
//   - word: The word pattern to search for in vocabulary entries.
//
// Returns:
//   - vocabulary_domain.SearchingResponse: A response containing the count of matching vocabulary entries and the vocabulary entries themselves.
//   - error: An error message if the search or conditions are not met.
func (v *vocabularyRepository) FetchByWordInBoth(ctx context.Context, word string) (vocabulary_domain.SearchingResponse, error) {
	vocabularySearchCh := make(chan vocabulary_domain.SearchingResponse, 1) // Channel to handle searching response
	wg.Add(1)                                                               // Add to wait group to manage goroutines
	go func() {
		defer wg.Done()
		data, found := vocabulariesSearchCache.Get(word)
		if found {
			vocabularySearchCh <- data
		}
	}()

	go func() {
		defer close(vocabularySearchCh) // Ensure channel is closed after processing
		wg.Wait()                       // Wait for all goroutines to finish
	}()

	vocabularySearchData := <-vocabularySearchCh // Receive data from channel
	if !internal.IsZeroValue(vocabularySearchData) {
		return vocabularySearchData, nil
	}

	// Get reference to the vocabulary collection from the database
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	// Define a regex pattern to search for the word case-insensitively
	regex := primitive.Regex{Pattern: word, Options: "i"}
	filter := bson.M{"word": bson.M{"$regex": regex}}

	var limit int64 = 10 // Limit the number of results to 10

	// Query the database to find vocabulary entries matching the filter
	cursor, err := collectionVocabulary.Find(ctx, filter, &options.FindOptions{Limit: &limit})
	if err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var vocabularies []vocabulary_domain.Vocabulary
	// Decode all results from the cursor into the vocabularies slice
	if err = cursor.All(ctx, &vocabularies); err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}

	// Count the total number of matching documents in the collection
	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}

	// Construct the response containing count and vocabulary entries
	vocabularyRes := vocabulary_domain.SearchingResponse{
		CountVocabularySearch: count,
		Vocabulary:            vocabularies,
	}

	// Cache the search result for future queries
	vocabulariesSearchCache.Set(word, vocabularyRes, 5*time.Minute)
	return vocabularyRes, nil
}

// FetchByLessonInBoth searches for vocabulary entries by lesson name pattern in an admin context.
// It first attempts to retrieve the data from the cache, and if not found, queries the database and caches the result.
// Parameters:
//   - ctx: The context for managing deadlines, cancellation signals, and other request-scoped values.
//   - lessonName: The lesson name pattern to search for in vocabulary entries.
//
// Returns:
//   - vocabulary_domain.SearchingResponse: A response containing the count of matching vocabulary entries and the vocabulary entries themselves.
//   - error: An error message if the search or conditions are not met.
func (v *vocabularyRepository) FetchByLessonInBoth(ctx context.Context, lessonName string) (vocabulary_domain.SearchingResponse, error) {
	vocabularySearchCh := make(chan vocabulary_domain.SearchingResponse, 1) // Channel to handle searching response
	wg.Add(1)                                                               // Add to wait group to manage goroutines
	go func() {
		defer wg.Done()
		data, found := vocabulariesSearchCache.Get(lessonName)
		if found {
			vocabularySearchCh <- data
		}
	}()

	go func() {
		defer close(vocabularySearchCh) // Ensure channel is closed after processing
		wg.Wait()                       // Wait for all goroutines to finish
	}()

	vocabularySearchData := <-vocabularySearchCh // Receive data from channel
	if !internal.IsZeroValue(vocabularySearchData) {
		return vocabularySearchData, nil
	}

	// Get reference to the vocabulary collection from the database
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	// Define a regex pattern to search for the lesson name case-insensitively
	regex := primitive.Regex{Pattern: lessonName, Options: "i"}
	filter := bson.M{"field_of_it": bson.M{"$regex": regex}}

	var limit int64 = 10 // Limit the number of results to 10

	// Query the database to find vocabulary entries matching the filter
	cursor, err := collectionVocabulary.Find(ctx, filter, &options.FindOptions{Limit: &limit})
	if err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var vocabularies []vocabulary_domain.Vocabulary
	// Decode all results from the cursor into the vocabularies slice
	if err := cursor.All(ctx, &vocabularies); err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}

	// Count the total number of matching documents in the collection
	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}

	// Construct the response containing count and vocabulary entries
	vocabularyRes := vocabulary_domain.SearchingResponse{
		CountVocabularySearch: count,
		Vocabulary:            vocabularies,
	}

	// Cache the search result for future queries
	vocabulariesSearchCache.Set(lessonName, vocabularyRes, 5*time.Minute)
	return vocabularyRes, nil
}

// FetchManyInBoth fetches multiple vocabulary entries from the database, with pagination support.
// It first attempts to retrieve cached data, and if not found, queries the database and caches the result.
// Parameters:
//   - ctx: The context for managing deadlines, cancellation signals, and other request-scoped values.
//   - page: The page number for pagination.
//
// Returns:
//   - vocabulary_domain.Response: The paginated response containing vocabulary entries.
//   - error: An error message if the fetch or conditions are not met.
func (v *vocabularyRepository) FetchManyInBoth(ctx context.Context, page string) (vocabulary_domain.Response, error) {
	errCh := make(chan error, 1)                               // Channel to handle errors
	vocabulariesCh := make(chan vocabulary_domain.Response, 1) // Channel to handle vocabulary response
	wg.Add(1)                                                  // Add a wait group to manage goroutines
	go func() {
		defer wg.Done()
		data, found := vocabularyResponseCache.Get(page)
		if found {
			vocabulariesCh <- data
		}
	}()

	go func() {
		defer close(vocabulariesCh) // Ensure vocabulariesCh is closed after processing
		wg.Wait()                   // Wait for all goroutines to finish
	}()

	vocabulariesData := <-vocabulariesCh // Receive data from the channel
	if !internal.IsZeroValue(vocabulariesData) {
		return vocabulariesData, nil
	}

	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return vocabulary_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	calCh := make(chan int64) // Channel to handle the total number of pages
	go func() {
		defer close(calCh)
		// Count total number of documents in the collection
		count, err := collectionVocabulary.CountDocuments(ctx, bson.D{})
		if err != nil {
			return
		}

		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)

		if cal2 != 0 {
			calCh <- cal1 + 1
		} else {
			calCh <- cal1
		}
	}()

	cursor, err := collectionVocabulary.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}(cursor, ctx)

	var vocabularies []vocabulary_domain.Vocabulary
	for cursor.Next(ctx) {
		var vocabulary vocabulary_domain.Vocabulary
		if err = cursor.Decode(&vocabulary); err != nil {
			return vocabulary_domain.Response{}, err
		}

		vocabularies = append(vocabularies, vocabulary)
	}

	cal := <-calCh // Receive total number of pages from the channel
	vocabularyRes := vocabulary_domain.Response{
		Page:        cal,
		CurrentPage: pageNumber,
		Vocabulary:  vocabularies,
	}

	vocabularyResponseCache.Set(page, vocabularyRes, 5*time.Minute) // Cache the response

	select {
	case err = <-errCh:
		return vocabulary_domain.Response{}, err
	default:
		return vocabularyRes, nil
	}
}

func (v *vocabularyRepository) CreateOneByNameUnitInAdmin(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	mu.Lock()
	defer mu.Unlock()

	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	isProcessing = true
	defer func() {
		isProcessing = false
	}()

	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)
	collectionLesson := v.database.Collection(v.collectionLesson)

	// Tìm unit dựa trên ID
	var unit unit_domain.Unit
	filterUnit := bson.M{"_id": vocabulary.UnitID}
	err := collectionUnit.FindOne(ctx, filterUnit).Decode(&unit)
	if err != nil {
		return err
	}

	filterLesson := bson.M{"_id": unit.LessonID}
	countLesson, err := collectionLesson.CountDocuments(ctx, filterLesson)
	if err != nil {
		return err
	}
	if countLesson == 0 {
		return errors.New("parent lesson not found")
	}

	filterUnit2 := bson.M{"_id": vocabulary.UnitID}
	countUnit, err := collectionUnit.CountDocuments(ctx, filterUnit2)
	if err != nil {
		return err
	}
	if countUnit == 0 {
		return errors.New("parent unit not found")
	}

	// Kiểm tra xem từ vựng đã tồn tại trong unit và bài học đó chưa
	filter := bson.M{"word": vocabulary.Word}
	countVocab, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if countVocab > 0 && countUnit > 0 {
		return errors.New("the vocabulary already exists in the lesson")
	}

	// Nếu không có lỗi, tạo bản ghi mới cho từ vựng
	_, err = collectionVocabulary.InsertOne(ctx, vocabulary)
	if err != nil {
		return err
	}

	return nil
}

// UpdateOneImageInAdmin updates the image URL of a vocabulary entry in the database within an admin context.
// It ensures thread safety by using a mutex and prevents concurrent processing with a flag.
// Parameters:
//   - ctx: The context for managing deadlines, cancellation signals, and other request-scoped values.
//   - vocabulary: The vocabulary entry containing the new image URL and the ID to be updated.
//
// Returns:
//   - *mongo.UpdateResult: The result of the update operation.
//   - error: An error message if the update fails or conditions are not met.
func (v *vocabularyRepository) UpdateOneImageInAdmin(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) (*mongo.UpdateResult, error) {
	// Lock to ensure no other goroutine is modifying at the same time
	mu.Lock()
	defer mu.Unlock()

	// Check if another goroutine is already processing
	if isProcessing {
		return nil, errors.New("another goroutine is already processing")
	}

	// Mark as processing
	isProcessing = true
	defer func() {
		isProcessing = false
	}()

	// Get the collection for vocabulary
	collection := v.database.Collection(v.collectionVocabulary)

	// Filter to find the vocabulary entry by ID
	filter := bson.D{{Key: "_id", Value: vocabulary.Id}}

	// Update operation to set the new image URL
	update := bson.M{
		"$set": bson.M{
			"image_url": vocabulary.LinkURL,
		},
	}

	// Perform the update operation
	data, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UpdateOneInAdmin updates various fields of a vocabulary entry in the database within an admin context.
// It ensures thread safety by using a mutex and prevents concurrent processing with a flag.
// Parameters:
//   - ctx: The context for managing deadlines, cancellation signals, and other request-scoped values.
//   - vocabulary: The vocabulary entry containing the new values and the ID to be updated.
//
// Returns:
//   - *mongo.UpdateResult: The result of the update operation.
//   - error: An error message if the update fails or conditions are not met.
func (v *vocabularyRepository) UpdateOneInAdmin(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) (*mongo.UpdateResult, error) {
	// Lock to ensure no other goroutine is modifying at the same time
	mu.Lock()
	defer mu.Unlock()

	// Check if another goroutine is already processing
	if isProcessing {
		return nil, errors.New("another goroutine is already processing")
	}

	// Mark as processing
	isProcessing = true
	defer func() {
		isProcessing = false
	}()

	// Get the collection for vocabulary
	collection := v.database.Collection(v.collectionVocabulary)

	// Filter to find the vocabulary entry by ID
	filter := bson.M{"_id": vocabulary.Id}

	// Update operation to set the new values
	update := bson.M{
		"$set": bson.M{
			"word":           vocabulary.Word,
			"part_of_speech": vocabulary.PartOfSpeech,
			"mean":           vocabulary.Mean,
			"pronunciation":  vocabulary.Pronunciation,
			"example_vie":    vocabulary.ExampleVie,
			"example_eng":    vocabulary.ExampleEng,
			"explain_vie":    vocabulary.ExplainVie,
			"explain_eng":    vocabulary.ExplainEng,
			"field_of_it":    vocabulary.FieldOfIT,
		},
	}

	// Perform the update operation
	data, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UpdateOneAudioInAdmin updates the audio link URL of a vocabulary entry in the database within an admin context.
// It ensures thread safety by using a mutex and prevents concurrent processing with a flag.
// Parameters:
//   - c: The context for managing deadlines, cancellation signals, and other request-scoped values.
//   - vocabulary: The vocabulary entry containing the new audio link URL and the ID to be updated.
//
// Returns:
//   - error: An error message if the update fails or conditions are not met.
func (v *vocabularyRepository) UpdateOneAudioInAdmin(c context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	// Lock to ensure no other goroutine is modifying at the same time
	mu.Lock()
	defer mu.Unlock()

	// Check if another goroutine is already processing
	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	// Mark as processing
	isProcessing = true
	defer func() {
		isProcessing = false
	}()

	// Get the collection for vocabulary
	collection := v.database.Collection(v.collectionVocabulary)

	// Filter to find the vocabulary entry by ID
	filter := bson.D{{Key: "_id", Value: vocabulary.Id}}

	// Update operation to set the new link URL
	update := bson.M{
		"$set": bson.M{
			"link_url": vocabulary.LinkURL,
		},
	}

	// Perform the update operation
	_, err := collection.UpdateOne(c, filter, &update)
	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyRepository) UpdateVocabularyProcess(ctx context.Context, vocabularyID string, process vocabulary_domain.VocabularyProcess) error {
	//TODO implement me
	panic("implement me")
}

func (v *vocabularyRepository) UpdateIsFavouriteInUser(ctx context.Context, vocabularyID string, isFavourite int) error {
	mu.Lock()
	defer mu.Unlock()

	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	isProcessing = true
	defer func() {
		isProcessing = false
	}()

	collection := v.database.Collection(v.collectionVocabulary)
	objID, err := primitive.ObjectIDFromHex(vocabularyID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.M{
		"$set": bson.M{
			"is_favourite": isFavourite,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// CreateOneInAdmin creates a new vocabulary entry in the database within an admin context.
// It uses a transaction to ensure atomicity of the operation, checking for the existence
// of the vocabulary and its parent unit before insertion.
// Parameters:
//   - ctx: The context for managing deadlines, cancelation signals, and other request-scoped values.
//   - vocabulary: The vocabulary entry to be inserted.
//
// Returns:
//   - error: An error message if the creation fails or conditions are not met.
func (v *vocabularyRepository) CreateOneInAdmin(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	// Lock to ensure no other goroutine is modifying at the same time
	mu.Lock()
	defer mu.Unlock()

	// Check if another goroutine is already processing
	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	// Mark as processing
	isProcessing = true
	defer func() {
		isProcessing = false
	}()

	// Start a new session for the transaction
	session, err := v.database.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Start the transaction
	err = session.StartTransaction()
	if err != nil {
		return err
	}

	// Get the collections for vocabulary and unit
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	// Filters to check for existing vocabulary and parent unit
	filter := bson.M{"word": vocabulary.Word, "unit_id": vocabulary.UnitID}
	filterReference := bson.M{"_id": vocabulary.UnitID}

	// Check if the parent unit exists
	countParent, err := collectionUnit.CountDocuments(ctx, filterReference)
	if err != nil {
		err := session.AbortTransaction(ctx)
		if err != nil {
			return err
		} // Abort the transaction if there's an error
		return err
	}

	// Check if the vocabulary already exists
	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		err := session.AbortTransaction(ctx)
		if err != nil {
			return err
		} // Abort the transaction if there's an error
		return errors.New("the vocabulary already exists")
	}
	if count > 0 {
		err := session.AbortTransaction(ctx)
		if err != nil {
			return err
		} // Abort the transaction if the vocabulary exists
		return errors.New("the word in unit already exists")
	}
	if countParent == 0 {
		err := session.AbortTransaction(ctx)
		if err != nil {
			return err
		} // Abort the transaction if the parent unit does not exist
		return errors.New("the parent unit does not exist")
	}

	// Insert the new vocabulary entry
	_, err = collectionVocabulary.InsertOne(ctx, vocabulary)
	if err != nil {
		err := session.AbortTransaction(ctx)
		if err != nil {
			return err
		} // Abort the transaction if there's an error during insertion
		return err
	}

	// Commit the transaction
	err = session.CommitTransaction(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteOneInAdmin deletes a vocabulary entry from the database in an admin context.
// It first checks if another goroutine is processing to avoid concurrent modifications.
// If the vocabulary entry has dependent marks, it will not be deleted.
// Parameters:
//   - ctx: The context for managing deadlines, cancelation signals, and other request-scoped values.
//   - vocabularyID: The ID of the vocabulary entry to be deleted.
//
// Returns:
//   - error: An error message if the deletion fails or conditions are not met.
func (v *vocabularyRepository) DeleteOneInAdmin(ctx context.Context, vocabularyID string) error {
	// Lock to ensure no other goroutine is modifying at the same time
	mu.Lock()
	defer mu.Unlock()

	// Check if another goroutine is already processing
	if isProcessing {
		return errors.New("another goroutine is already processing")
	}

	// Mark as processing
	isProcessing = true
	defer func() {
		isProcessing = false
	}()

	// Get the collections for vocabulary and mark
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionMark := v.database.Collection(v.collectionMark)

	// Convert vocabularyID to ObjectID
	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return err
	}

	// Filters for the vocabulary and its child marks
	filter := bson.M{
		"_id": objID,
	}

	filterChild := bson.M{
		"vocabulary_id": objID,
	}

	// Check if there are child marks associated with the vocabulary entry
	countChildMark, err := collectionMark.CountDocuments(ctx, filterChild)
	if err != nil {
		return err
	}
	if countChildMark > 0 {
		return errors.New("lesson cannot remove")
	}

	// Check if the vocabulary entry exists
	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("the lesson is removed")
	}

	// Delete the vocabulary entry
	_, err = collectionVocabulary.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return err
}
