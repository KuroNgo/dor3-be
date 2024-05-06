package vocabulary_repository

import (
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/internal"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"sync"
	"time"
)

type vocabularyRepository struct {
	database             *mongo.Database
	collectionVocabulary string
	collectionMark       string
	collectionUnit       string

	vocabularyManyCache    map[string]vocabulary_domain.Response
	vocabularyOneCache     map[string][]vocabulary_domain.Vocabulary
	vocabularyCacheExpires map[string]time.Time
	cacheMutex             sync.RWMutex
}

func (v *vocabularyRepository) GetVocabularyById(ctx context.Context, id string) (vocabulary_domain.Vocabulary, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	idVocabulary, err := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": idVocabulary}
	if err != nil {
		return vocabulary_domain.Vocabulary{}, err
	}

	var vocabulary vocabulary_domain.Vocabulary
	err = collectionVocabulary.FindOne(ctx, filter).Decode(&vocabulary)
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return vocabulary_domain.Vocabulary{}, err
	}

	return vocabulary, nil
}

func NewVocabularyRepository(db *mongo.Database, collectionVocabulary string, collectionMark string, collectionUnit string) vocabulary_domain.IVocabularyRepository {
	return &vocabularyRepository{
		database:             db,
		collectionVocabulary: collectionVocabulary,
		collectionMark:       collectionMark,
		collectionUnit:       collectionUnit,

		vocabularyManyCache:    make(map[string]vocabulary_domain.Response),
		vocabularyOneCache:     make(map[string][]vocabulary_domain.Vocabulary),
		vocabularyCacheExpires: make(map[string]time.Time),
	}
}

func (v *vocabularyRepository) FindUnitIDByUnitLevel(ctx context.Context, unitLevel int) (primitive.ObjectID, error) {
	collectionUnit := v.database.Collection(v.collectionUnit)

	// Tìm kiếm unit có cùng level
	filter := bson.M{"level": unitLevel}
	var existingUnit struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	err := collectionUnit.FindOne(ctx, filter).Decode(&existingUnit)
	if err == nil {
		return existingUnit.Id, nil
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return primitive.NilObjectID, err
	}

	newUnit := unit_domain.Unit{
		ID:         primitive.NewObjectID(),
		Name:       "Unit " + strconv.Itoa(unitLevel),
		Level:      unitLevel,
		IsComplete: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = collectionUnit.InsertOne(ctx, newUnit)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return newUnit.ID, nil
}

func (v *vocabularyRepository) FindVocabularyIDByVocabularyConfig(ctx context.Context, word string) (primitive.ObjectID, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	filter := bson.M{"word_for_config": word}
	var data struct {
		Id primitive.ObjectID `bson:"_id"`
	}

	err := collectionVocabulary.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return data.Id, nil
}

func (v *vocabularyRepository) GetLatestVocabulary(ctx context.Context) ([]string, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	var vocabularies []string
	yesterday := time.Now().Add(-24 * time.Hour)
	filter := bson.M{"created_at": bson.M{"$gt": yesterday}}

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

	for cursor.Next(ctx) {
		var result bson.M
		if err = cursor.Decode(&result); err != nil {
			return nil, err
		}
		word, ok := result["word"].(string)
		if !ok {
			return nil, errors.New("failed to parse word from result")
		}
		vocabularies = append(vocabularies, word)
	}

	return vocabularies, nil
}

func (v *vocabularyRepository) GetAllVocabulary(ctx context.Context) ([]string, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	var vocabularies []string

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

	for cursor.Next(ctx) {
		var result bson.M
		if err = cursor.Decode(&result); err != nil {
			return nil, err
		}
		word, ok := result["word"].(string)
		if !ok {
			return nil, errors.New("failed to parse word from result")
		}
		vocabularies = append(vocabularies, word)
	}

	return vocabularies, nil
}

func (v *vocabularyRepository) CreateOneByNameUnit(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	filter := bson.M{"word": vocabulary.Word, "unit_id": vocabulary.UnitID}

	filterParent := bson.M{"_id": vocabulary.UnitID}
	countParent, err := collectionUnit.CountDocuments(ctx, filterParent)
	if err != nil {
		return err
	}
	if countParent == 0 {
		return errors.New("parent unit not found")
	}

	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the vocabulary already exists in the lesson")
	}

	_, err = collectionVocabulary.InsertOne(ctx, vocabulary)
	if err != nil {
		return err
	}
	return nil
}

func (v *vocabularyRepository) FetchByIdUnit(ctx context.Context, idUnit string) (vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	idUnit2, err := primitive.ObjectIDFromHex(idUnit)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}

	filter := bson.M{"unit_id": idUnit2}

	cursor, err := collectionVocabulary.Find(ctx, filter)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var vocabularies []vocabulary_domain.Vocabulary

	for cursor.Next(ctx) {
		var vocabulary vocabulary_domain.Vocabulary
		if err = cursor.Decode(&vocabulary); err != nil {
			return vocabulary_domain.Response{}, err
		}
		vocabulary.UnitID = idUnit2
		vocabularies = append(vocabularies, vocabulary)
	}

	response := vocabulary_domain.Response{
		Vocabulary: vocabularies,
	}
	return response, nil
}

func (v *vocabularyRepository) FetchByWord(ctx context.Context, word string) (vocabulary_domain.SearchingResponse, error) {
	v.cacheMutex.RLock()
	cachedData, found := v.vocabularyOneCache[word]
	v.cacheMutex.RUnlock()

	if found {
		vocabularyRes := vocabulary_domain.SearchingResponse{
			CountVocabularySearch: int64(len(cachedData)),
			Vocabulary:            cachedData,
		}
		return vocabularyRes, nil
	}

	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	regex := primitive.Regex{Pattern: word, Options: "i"}
	filter := bson.M{"word": bson.M{"$regex": regex}}

	var limit int64 = 10

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
	if err := cursor.All(ctx, &vocabularies); err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}

	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}

	vocabularyRes := vocabulary_domain.SearchingResponse{
		CountVocabularySearch: count,
		Vocabulary:            vocabularies,
	}

	v.cacheMutex.Lock()
	v.vocabularyOneCache[word] = vocabularies
	v.vocabularyCacheExpires[word] = time.Now().Add(5 * time.Minute)
	v.cacheMutex.Unlock()

	return vocabularyRes, nil
}

func (v *vocabularyRepository) FetchByLesson(ctx context.Context, lessonName string) (vocabulary_domain.SearchingResponse, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)

	regex := primitive.Regex{Pattern: lessonName, Options: "i"}
	filter := bson.M{"field_of_it": bson.M{"$regex": regex}}

	var limit int64 = 10

	cursor, err := collectionVocabulary.Find(ctx, filter, &options.FindOptions{Limit: &limit})
	if err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}
	defer cursor.Close(ctx)

	var vocabularies []vocabulary_domain.Vocabulary
	if err := cursor.All(ctx, &vocabularies); err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}

	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return vocabulary_domain.SearchingResponse{}, err
	}

	vocabularyRes := vocabulary_domain.SearchingResponse{
		CountVocabularySearch: count,
		Vocabulary:            vocabularies,
	}

	return vocabularyRes, nil
}

func (v *vocabularyRepository) FetchMany(ctx context.Context, page string) (vocabulary_domain.Response, error) {
	// Kiểm tra cache trước khi truy vấn cơ sở dữ liệu
	v.cacheMutex.RLock()
	cachedData, found := v.vocabularyManyCache[page]
	v.cacheMutex.RUnlock()

	if found {
		return cachedData, nil
	}

	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return vocabulary_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	calCh := make(chan int64)
	go func() {
		defer close(calCh)
		// Đếm tổng số lượng tài liệu trong collection
		count, err := collectionVocabulary.CountDocuments(ctx, bson.D{})
		if err != nil {
			return
		}

		cal1 := count / int64(perPage)
		cal2 := count % int64(perPage)

		if cal2 != 0 {
			calCh <- cal1
		}
	}()

	cursor, err := collectionVocabulary.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}
	defer cursor.Close(ctx)

	var vocabularies []vocabulary_domain.Vocabulary

	for cursor.Next(ctx) {
		var vocabulary vocabulary_domain.Vocabulary

		if err := cursor.Decode(&vocabulary); err != nil {
			return vocabulary_domain.Response{}, err
		}

		var unit unit_domain.Unit
		if err := collectionUnit.FindOne(ctx, bson.M{"_id": vocabulary.UnitID}).Decode(&unit); err != nil {
			return vocabulary_domain.Response{}, err
		}

		vocabulary.UnitID = unit.ID
		vocabularies = append(vocabularies, vocabulary)
	}

	cal := <-calCh
	vocabularyRes := vocabulary_domain.Response{
		Page:       cal,
		Vocabulary: vocabularies,
	}

	v.cacheMutex.Lock()
	v.vocabularyManyCache[page] = vocabularyRes
	v.vocabularyCacheExpires[page] = time.Now().Add(5 * time.Minute) // Ví dụ: hết hạn sau 5 phút
	v.cacheMutex.Unlock()

	return vocabularyRes, nil
}

func (v *vocabularyRepository) UpdateOne(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) (*mongo.UpdateResult, error) {
	collection := v.database.Collection(v.collectionVocabulary)

	filter := bson.M{"_id": vocabulary.Id}
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

	data, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (v *vocabularyRepository) UpdateOneAudio(c context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	collection := v.database.Collection(v.collectionVocabulary)

	filter := bson.D{{Key: "_id", Value: vocabulary.Id}}
	update := bson.M{
		"$set": bson.M{
			"link_url": vocabulary.LinkURL,
		},
	}

	_, err := collection.UpdateOne(c, filter, &update)
	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyRepository) UpdateIsFavourite(ctx context.Context, vocabularyID string, isFavourite int) error {
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

func (v *vocabularyRepository) CreateOne(ctx context.Context, vocabulary *vocabulary_domain.Vocabulary) error {
	session, err := v.database.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Bắt đầu transaction
	err = session.StartTransaction()
	if err != nil {
		return err
	}

	// Thực hiện các thao tác dữ liệu trong transaction
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	filter := bson.M{"word": vocabulary.Word, "unit_id": vocabulary.UnitID}
	filterReference := bson.M{"_id": vocabulary.UnitID}

	countParent, err := collectionUnit.CountDocuments(ctx, filterReference)
	if err != nil {
		err := session.AbortTransaction(ctx)
		if err != nil {
			return err
		}
		return err
	}

	// Kiểm tra và thêm từ vựng vào cơ sở dữ liệu
	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		err := session.AbortTransaction(ctx)
		if err != nil {
			return errors.New("the vocabulary already exist")
		}
		return err
	}
	if count > 0 {
		err := session.AbortTransaction(ctx)
		if err != nil {
			return err
		}
		return errors.New("the word in unit already exists")
	}
	if countParent == 0 {
		err := session.AbortTransaction(ctx)
		if err != nil {
			return err
		}
		return errors.New("the parent unit does not exist")
	}

	_, err = collectionVocabulary.InsertOne(ctx, vocabulary)
	if err != nil {
		err := session.AbortTransaction(ctx)
		if err != nil {
			return err
		}
		return err
	}

	// Kết thúc transaction
	err = session.CommitTransaction(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (v *vocabularyRepository) UpsertOne(ctx context.Context, id string, vocabulary *vocabulary_domain.Vocabulary) (vocabulary_domain.Response, error) {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionUnit := v.database.Collection(v.collectionUnit)

	filterReference := bson.M{"_id": vocabulary.UnitID}
	count, err := collectionUnit.CountDocuments(ctx, filterReference)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}

	if count == 0 {
		return vocabulary_domain.Response{}, errors.New("the lesson ID do not exist")
	}

	doc, err := internal.ToDoc(vocabulary)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return vocabulary_domain.Response{}, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "_id", Value: idHex}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := collectionVocabulary.FindOneAndUpdate(ctx, query, update, opts)

	var updatePost vocabulary_domain.Response
	if err := res.Decode(&updatePost); err != nil {
		return vocabulary_domain.Response{}, err
	}

	return updatePost, nil
}

func (v *vocabularyRepository) DeleteOne(ctx context.Context, vocabularyID string) error {
	collectionVocabulary := v.database.Collection(v.collectionVocabulary)
	collectionMark := v.database.Collection(v.collectionMark)

	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}

	filterChild := bson.M{
		"vocabulary_id": objID,
	}

	countChildMark, err := collectionMark.CountDocuments(ctx, filterChild)
	if err != nil {
		return err
	}
	if countChildMark > 0 {
		return errors.New(`lesson cannot remove`)
	}

	count, err := collectionVocabulary.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`the lesson is removed`)
	}

	_, err = collectionVocabulary.DeleteOne(ctx, filter)
	return err
}

func (v *vocabularyRepository) DeleteMany(ctx context.Context, vocabularyID ...string) error {
	//TODO implement me
	panic("implement me")
}
