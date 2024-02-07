package router

import (
	user_router "clean-architecture/api/router/user"
	"clean-architecture/bootstrap"
	"clean-architecture/infrastructor/mongo"
	"github.com/gin-gonic/gin"
	"time"
)

func Setup(env *bootstrap.Database, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	publicRouter := gin.Group("")

	// All Public APIs

	// user method
	user_router.NewLoginRouter(env, timeout, db, publicRouter)

	// Middleware to verify accessToken

	// All Protected APIs

	// All Private APIs
}
