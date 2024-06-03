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

	//var memStats runtime.MemStats
	//runtime.ReadMemStats(&memStats)
	//fmt.Printf("Alloc = %v MiB\n", bToMb(memStats.Alloc))
	//fmt.Printf("TotalAlloc = %v MiB\n", bToMb(memStats.TotalAlloc))
	//fmt.Printf("Sys = %v MiB\n", bToMb(memStats.Sys))
	//fmt.Printf("NumGC = %v\n", memStats.NumGC)

	err := _gin.Run(env.ServerAddress)
	if err != nil {
		return
	}

}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
