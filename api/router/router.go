package router

import (
	user_controller "clean-architecture/api/controller/user"
	"clean-architecture/api/middleware"
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
	user_router.GoogleAuthRouter(env, timeout, db, publicRouter)
	publicRouter.GET("/logout", user_controller.LogoutUser)

	// Middleware
	publicRouter.OPTIONS("/*path", middleware.OptionMessage)
	publicRouter.Use(middleware.CORSPublic())
	// All Protected APIs

	// All Private APIs
}
