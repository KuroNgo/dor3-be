package mongo

import (
	"clean-architecture/bootstrap"
	"context"
	"fmt"
	mongo_driven "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func NewMongoDatabase(env *bootstrap.Database) *mongo_driven.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//dbHost := env.DBHost
	//dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPassword

	mongodbURI := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.ykpyhgp.mongodb.net/?authMechanism=SCRAM-SHA-1", dbUser, dbPass)

	//if dbUser == "" || dbPass == "" {
	//	mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	//}

	mongoCon := options.Client().ApplyURI(mongodbURI)
	client, err := mongo_driven.Connect(ctx, mongoCon)
	if err != nil {
		log.Fatal("error while connecting with mongo", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}

	session, err := client.StartSession()
	if err != nil {
		log.Fatal("err to start sessions", err)
	}
	defer session.EndSession(context.Background())

	err = session.StartTransaction()
	if err != nil {
		log.Fatal("err to start sessions", err)
	}
	defer func() {
		if err != nil {
			// Rollback giao dịch
			err := session.AbortTransaction(context.Background())
			if err != nil {
				return
			}
			return
		}
		// Commit giao dịch
		err := session.CommitTransaction(context.Background())
		if err != nil {
			return
		}
	}()
	return client
}

// test
//func NewMongoDatabase(env *bootstrap.Database) *mongo_driven.Client {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	mongodbURI := fmt.Sprintf("mongodb://localhost:27017/")
//
//	mongoCon := options.Client().ApplyURI(mongodbURI)
//	client, err := mongo_driven.Connect(ctx, mongoCon)
//	if err != nil {
//		log.Fatal("error while connecting with mongo", err)
//	}
//
//	err = client.Ping(ctx, readpref.Primary())
//	if err != nil {
//		log.Fatal("error while trying to ping mongo", err)
//	}
//
//	session, err := client.StartSession()
//	if err != nil {
//		log.Fatal("err to start sessions", err)
//	}
//	defer session.EndSession(context.Background())
//
//	err = session.StartTransaction()
//	if err != nil {
//		log.Fatal("err to start sessions", err)
//	}
//	defer func() {
//		if err != nil {
//			// Rollback giao dịch
//			err := session.AbortTransaction(context.Background())
//			if err != nil {
//				return
//			}
//			return
//		}
//		// Commit giao dịch
//		err := session.CommitTransaction(context.Background())
//		if err != nil {
//			return
//		}
//	}()
//	return client
//}

func CloseMongoDBConnection(client *mongo_driven.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
