package router

import (
	"clean-architecture/bootstrap"
	"clean-architecture/infrastructor/mongo"
	"github.com/gin-gonic/gin"
	"time"
)

func Setup(env *bootstrap.Database, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	//publicRouter := gin.Group("")
	// All Public APIs

	//protectedRouter := gin.Group("")
	// Middleware to verify AccessToken

}
