package jade_repository

import (
	jade_domain "clean-architecture/domain/jade"
	blockchain2 "clean-architecture/internal/blockchain"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type jadeRepository struct {
	database       *mongo.Database
	collectionJade string
}

func NewJadeRepository(db *mongo.Database, collectionJade string) jade_domain.IJadeRepository {
	return &jadeRepository{
		database:       db,
		collectionJade: collectionJade,
	}
}

func (j *jadeRepository) FetchJadeInUser(ctx context.Context, userID primitive.ObjectID) (jade_domain.JadeBlockchain, error) {
	collectionJade := j.database.Collection(j.collectionJade)

	filter := bson.M{"user_id": userID}
	var jade jade_domain.Jade
	err := collectionJade.FindOne(ctx, filter).Decode(&jade)
	if err != nil {
		return jade_domain.JadeBlockchain{}, err
	}

	var block *blockchain2.Block
	block = jade.Amount

	// Accessing data of each block
	jadeBlockchain := jade_domain.JadeBlockchain{
		ID:     jade.ID,
		UserID: jade.UserID,
		Amount: block.Data,
	}

	return jadeBlockchain, nil
}

func (j *jadeRepository) Rank(ctx context.Context) ([]jade_domain.JadeBlockchain, error) {
	collectionJade := j.database.Collection(j.collectionJade)

	findOptions := options.Find().SetSort(bson.M{"amount": 1})
	cursor, err := collectionJade.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}

	var jades []jade_domain.JadeBlockchain
	jades = make([]jade_domain.JadeBlockchain, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var jade jade_domain.Jade
		err = cursor.Decode(&jade)
		if err != nil {
			return nil, err
		}

		jadeBlc := jade_domain.JadeBlockchain{
			ID:     jade.ID,
			UserID: jade.UserID,
			Amount: jade.Amount.Data,
		}

		jades = append(jades, jadeBlc)
	}

	return jades, nil
}

func (j *jadeRepository) CreateCurrencyInUser(ctx context.Context, userID primitive.ObjectID, data int64) error {
	collectionJade := j.database.Collection(j.collectionJade)

	bc := blockchain2.NewBlockchain()
	bc.AddBlock(data)
	newBlock := bc.Blocks[len(bc.Blocks)-1]

	jade := jade_domain.Jade{
		ID:     primitive.NewObjectID(),
		UserID: userID,
		Amount: newBlock,
	}

	_, err := collectionJade.InsertOne(ctx, &jade)
	if err != nil {
		return err
	}

	return nil
}

func (j *jadeRepository) UpdateCurrencyInUser(ctx context.Context, userID primitive.ObjectID, data int64) error {
	collectionJade := j.database.Collection(j.collectionJade)

	filter := bson.M{"user_id": userID}
	bc := blockchain2.NewBlockchain()
	bc.AddBlock(data)
	newBlock := bc.Blocks[len(bc.Blocks)-1]
	update := bson.M{
		"$set": bson.M{
			"amount": newBlock,
		},
	}

	_, err := collectionJade.UpdateOne(ctx, filter, &update)
	if err != nil {
		return err
	}

	return nil
}
