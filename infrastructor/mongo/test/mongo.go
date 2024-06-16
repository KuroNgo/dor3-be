package mongo

import (
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// Database interface for interacting with a MongoDB database
//
//go:generate mockery --name Database
type Database interface {
	Collection(string) Collection
	Client() Client
}

// Collection interface for interacting with a MongoDB collection
//
//go:generate mockery --name Collection
type Collection interface {
	FindOne(context.Context, interface{}) SingleResult
	InsertOne(context.Context, interface{}) (interface{}, error)
	InsertMany(context.Context, []interface{}) ([]interface{}, error)
	DeleteOne(context.Context, interface{}) (int64, error)
	Find(context.Context, interface{}, ...*options.FindOptions) (Cursor, error)
	CountDocuments(context.Context, interface{}, ...*options.CountOptions) (int64, error)
	Aggregate(context.Context, interface{}) (Cursor, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

// SingleResult interface for handling a single result from a query
//
//go:generate mockery --name SingleResult
type SingleResult interface {
	Decode(interface{}) error
}

// Cursor interface for handling multiple results from a query
//
//go:generate mockery --name Cursor
type Cursor interface {
	Close(context.Context) error
	Next(context.Context) bool
	Decode(interface{}) error
	All(context.Context, interface{}) error
}

// Client interface for interacting with a MongoDB client
//
//go:generate mockery --name Client
type Client interface {
	Database(string) Database
	Connect(context.Context) error
	Disconnect(context.Context) error
	StartSession() (mongo.Session, error)
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
	Ping(context.Context) error
}

// Define mock structs
type MockCollection struct {
	mock.Mock
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Clear() {
	m.Called()
}
