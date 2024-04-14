package main

import (
	"clean-architecture/api/router"
	"clean-architecture/infrastructor/mongo"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {

	app := mongo.App()

	env := app.Env

	db := app.MongoDB.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	_gin := gin.Default()

	router.SetUp(env, timeout, db, _gin)

	err := _gin.Run(env.ServerAddress)
	if err != nil {
		return
	}

}
