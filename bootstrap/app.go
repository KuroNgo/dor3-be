package bootstrap

import (
	mongo_driven "go.mongodb.org/mongo-driver/mongo"
)

//type Application struct {
//	Env   *Database
//	Mongo mongo.Client
//}

type Application struct {
	Env     *Database
	MongoDB *mongo_driven.Client
}

//func App() Application {
//	app := &Application{}
//	app.Env = NewEnv()
//	app.Mongo = NewMongoDatabase(app.Env)
//	return *app
//}

func App() *Application {
	app := &Application{}
	app.Env = NewEnv()
	app.MongoDB = NewMongoDatabase(app.Env)
	return app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.MongoDB)
}
