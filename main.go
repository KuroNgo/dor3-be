package main

import (
	"clean-architecture/api/router"
	"clean-architecture/bootstrap"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	_gin := gin.Default()

	router.SetUp(env, timeout, db, _gin)
	err := _gin.Run(env.ServerAddress)
	if err != nil {
		return
	}
}
