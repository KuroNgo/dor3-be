package mongo

import (
	"clean-architecture/bootstrap"
	mongo_driven "go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Env     *bootstrap.Database
	MongoDB *mongo_driven.Client
}

func App() *Application {
	app := &Application{}
	app.Env = bootstrap.NewEnv()
	app.MongoDB = NewMongoDatabase(app.Env)
	return app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.MongoDB)
}
