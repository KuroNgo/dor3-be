package router

import (
	user_controller "clean-architecture/api/controller/user"
	"clean-architecture/api/middleware"
	quiz_route "clean-architecture/api/router/quiz"
	user_router "clean-architecture/api/router/user"
	"clean-architecture/bootstrap"
	"clean-architecture/infrastructor/mongo"
	"github.com/gin-gonic/gin"
	"time"
)

func Setup(env *bootstrap.Database, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	publicRouter := gin.Group("")
	privateRouter := gin.Group("/admin")

	// Middleware
	publicRouter.Use(
		middleware.CORSPublic(),
		middleware.RateLimiter(),
	)

	privateRouter.Use(
		middleware.CORSForDev(),
		middleware.RateLimiter(),
		//middleware.DeserializeUser(),
	)

	// This is a CORS method for check IP valid
	publicRouter.OPTIONS("/*path", middleware.OptionMessage)

	// All Public APIs
	// user method
	user_router.GoogleAuthRouter(env, timeout, db, publicRouter)
	user_router.RefreshTokenRouter(env, timeout, db, publicRouter)

	// quiz method
	quiz_route.QuizRouter(env, timeout, db, publicRouter)
	publicRouter.GET("/logout", middleware.DeserializeUser(), user_controller.LogoutUser)

	// private router
	quiz_route.AdminQuizRouter(env, timeout, db, privateRouter)

}
