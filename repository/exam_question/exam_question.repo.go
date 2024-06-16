package exam_question_repository

import (
	exam_question_domain "clean-architecture/domain/exam_question"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type examQuestionRepository struct {
	database             *mongo.Database
	collectionQuestion   string
	collectionExam       string
	collectionVocabulary string
}

func NewExamQuestionRepository(db *mongo.Database, collectionQuestion string, collectionExam string, collectionVocabulary string) exam_question_domain.IExamQuestionRepository {
	return &examQuestionRepository{
		database:             db,
		collectionQuestion:   collectionQuestion,
		collectionExam:       collectionExam,
		collectionVocabulary: collectionVocabulary,
	}
}

func (e *examQuestionRepository) FetchManyInAdmin(ctx context.Context, page string) (exam_question_domain.Response, error) {
	collectionQuestion := e.database.Collection(e.collectionQuestion)
	collectVocabulary := e.database.Collection(e.collectionVocabulary)

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return exam_question_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	count, err := collectionQuestion.CountDocuments(ctx, bson.D{})
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	cal1 := count / int64(perPage)
	cal2 := count % int64(perPage)
	var cal int64 = 0
	if cal2 != 0 {
		cal = cal1 + 1
	}

	cursor, err := collectionQuestion.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return exam_question_domain.Response{}, err
	}

	var questions []exam_question_domain.ExamQuestionResponse
	for cursor.Next(ctx) {
		var question exam_question_domain.ExamQuestion
		if err := cursor.Decode(&question); err != nil {
			fmt.Println("Error decoding question:", err)
			return exam_question_domain.Response{}, err
		}

		var question2 exam_question_domain.ExamQuestionResponse
		if err := cursor.Decode(&question2); err != nil {
			fmt.Println("Error decoding question:", err)
			return exam_question_domain.Response{}, err
		}

		var vocabulary vocabulary_domain.Vocabulary
		filterVocabulary := bson.M{"_id": question.VocabularyID}
		_ = collectVocabulary.FindOne(ctx, filterVocabulary).Decode(&vocabulary)

		question2.Vocabulary = vocabulary
		questions = append(questions, question2)
	}
	questionsRes := exam_question_domain.Response{
		Page:                 cal,
		ExamQuestionResponse: questions,
	}
	return questionsRes, nil
}

func (e *examQuestionRepository) FetchOneByExamIDInAdmin(ctx context.Context, examID string) (exam_question_domain.ExamQuestionResponse, error) {
	collectionQuestion := e.database.Collection(e.collectionQuestion)
	collectionVocabulary := e.database.Collection(e.collectionVocabulary)

	idExam, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return exam_question_domain.ExamQuestionResponse{}, err
	}

	var examQuestion exam_question_domain.ExamQuestion
	filter := bson.M{"exam_id": idExam}
	err = collectionQuestion.FindOne(ctx, filter).Decode(&examQuestion)
	if err != nil {
		return exam_question_domain.ExamQuestionResponse{}, err
	}

	var vocabulary vocabulary_domain.Vocabulary
	filterVocabulary := bson.M{"_id": examQuestion.VocabularyID}
	_ = collectionVocabulary.FindOne(ctx, filterVocabulary).Decode(&vocabulary)

	examQuestionRes := exam_question_domain.ExamQuestionResponse{
		ID:            examQuestion.ID,
		ExamID:        examQuestion.ExamID,
		Vocabulary:    vocabulary,
		Content:       examQuestion.Content,
		Type:          examQuestion.Type,
		Level:         examQuestion.Level,
		Options:       examQuestion.Options,
		CorrectAnswer: examQuestion.CorrectAnswer,
		CreatedAt:     examQuestion.CreatedAt,
		UpdateAt:      examQuestion.UpdateAt,
		WhoUpdate:     examQuestion.WhoUpdate,
	}

	return examQuestionRes, nil
}

func (e *examQuestionRepository) FetchQuestionByIDInAdmin(ctx context.Context, id string) (exam_question_domain.ExamQuestion, error) {
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	idQuestion, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return exam_question_domain.ExamQuestion{}, err
	}

	var examQuestion exam_question_domain.ExamQuestion
	filter := bson.M{"_id": idQuestion}
	err = collectionQuestion.FindOne(ctx, filter).Decode(&examQuestion)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return exam_question_domain.ExamQuestion{}, errors.New("question not found")
		}
		return exam_question_domain.ExamQuestion{}, err
	}

	return examQuestion, nil
}

func (e *examQuestionRepository) FetchManyByExamIDInAdmin(ctx context.Context, examID string, page string) (exam_question_domain.Response, error) {
	collectionQuestion := e.database.Collection(e.collectionQuestion)
	collectVocabulary := e.database.Collection(e.collectionVocabulary)

	idExam, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		fmt.Println("Error converting examID to ObjectID:", err)
		return exam_question_domain.Response{}, err
	}

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println("Error converting page to int:", err)
		return exam_question_domain.Response{}, errors.New("invalid page number")
	}
	perPage := 7
	skip := (pageNumber - 1) * perPage
	findOptions := options.Find().SetLimit(int64(perPage)).SetSkip(int64(skip))

	filter := bson.M{"exam_id": idExam}
	cursor, err := collectionQuestion.Find(ctx, filter, findOptions)
	if err != nil {
		fmt.Println("Error finding documents in collectionQuestion:", err)
		return exam_question_domain.Response{}, err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Println("Error closing cursor:", err)
		}
	}()

	var count int64
	count, err = collectionQuestion.CountDocuments(ctx, bson.M{"exam_id": idExam})
	if err != nil {
		fmt.Println("Error counting documents in collectionQuestion:", err)
		return exam_question_domain.Response{}, err
	}

	var questions []exam_question_domain.ExamQuestionResponse

	for cursor.Next(ctx) {
		var question exam_question_domain.ExamQuestionResponse
		if err := cursor.Decode(&question); err != nil {
			fmt.Println("Error decoding question:", err)
			return exam_question_domain.Response{}, err
		}

		var question2 exam_question_domain.ExamQuestion
		if err := cursor.Decode(&question2); err != nil {
			fmt.Println("Error decoding question:", err)
			return exam_question_domain.Response{}, err
		}

		var vocabulary vocabulary_domain.Vocabulary
		filterVocabulary := bson.M{"_id": question2.VocabularyID}
		_ = collectVocabulary.FindOne(ctx, filterVocabulary).Decode(&vocabulary)

		question.Vocabulary = vocabulary
		questions = append(questions, question)
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("Cursor encountered an error:", err)
		return exam_question_domain.Response{}, err
	}

	var totalPages int64
	if count%int64(perPage) == 0 {
		totalPages = count / int64(perPage)
	} else {
		totalPages = count/int64(perPage) + 1
	}

	statisticsCh := make(chan exam_question_domain.Statistics)
	go func() {
		statistics, _ := e.Statistics(ctx)
		statisticsCh <- statistics
	}()
	statistics := <-statisticsCh

	questionsRes := exam_question_domain.Response{
		Statistics:           statistics,
		Page:                 totalPages,
		CurrentPage:          pageNumber,
		ExamQuestionResponse: questions,
	}

	return questionsRes, nil
}

func (e *examQuestionRepository) UpdateOneInAdmin(ctx context.Context, examQuestion *exam_question_domain.ExamQuestion) (*mongo.UpdateResult, error) {
	collection := e.database.Collection(e.collectionQuestion)

	filter := bson.D{{Key: "_id", Value: examQuestion.ID}}
	update := bson.M{
		"$set": bson.M{
			"content":        examQuestion.Content,
			"type":           examQuestion.Type,
			"level":          examQuestion.Level,
			"options":        examQuestion.Options,
			"correct_answer": examQuestion.CorrectAnswer,
			"update_at":      examQuestion.UpdateAt,
			"who_update":     examQuestion.WhoUpdate,
		},
	}

	data, err := collection.UpdateOne(ctx, filter, &update)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *examQuestionRepository) CreateOneInAdmin(ctx context.Context, examQuestion *exam_question_domain.ExamQuestion) error {
	collectionQuestion := e.database.Collection(e.collectionQuestion)
	collectionExam := e.database.Collection(e.collectionExam)
	collectionVocabulary := e.database.Collection(e.collectionVocabulary)

	filterExamID := bson.M{"_id": examQuestion.ExamID}
	countExamID, err := collectionExam.CountDocuments(ctx, filterExamID)
	if err != nil {
		return err
	}
	if countExamID == 0 {
		return errors.New("the examID does not exist")
	}

	filterVocabularyID := bson.M{"_id": examQuestion.VocabularyID}
	countVocabularyID, err := collectionVocabulary.CountDocuments(ctx, filterVocabularyID)
	if err != nil {
		return err
	}
	if countVocabularyID == 0 {
		return errors.New("the vocabularyID does not exist")
	}

	filterParent := bson.M{"exam_id": examQuestion.ExamID}
	count, err := collectionQuestion.CountDocuments(ctx, filterParent)
	if err != nil {
		return err
	}
	if count >= 10 {
		return errors.New("the question id is not added in one exam")
	}

	// Thêm câu hỏi vào cơ sở dữ liệu nếu không có lỗi
	_, err = collectionQuestion.InsertOne(ctx, examQuestion)
	if err != nil {
		return err
	}
	return nil
}

func (e *examQuestionRepository) DeleteOneInAdmin(ctx context.Context, examID string) error {
	collectionQuestion := e.database.Collection(e.collectionQuestion)

	objID, err := primitive.ObjectIDFromHex(examID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	count, err := collectionQuestion.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`exam answer is removed`)
	}

	_, err = collectionQuestion.DeleteOne(ctx, filter)
	return err
}

func (e *examQuestionRepository) Statistics(ctx context.Context) (exam_question_domain.Statistics, error) {
	collectionExamQuestion := e.database.Collection(e.collectionQuestion)

	count, err := collectionExamQuestion.CountDocuments(ctx, bson.D{})
	if err != nil {
		return exam_question_domain.Statistics{}, err
	}

	statistics := exam_question_domain.Statistics{
		Count: count,
	}
	return statistics, nil
}
