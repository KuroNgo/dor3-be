package router

import (
	user_controller "clean-architecture/api/controller/user"
	"clean-architecture/api/middleware"
	audio_route "clean-architecture/api/router/audio"
	course_route "clean-architecture/api/router/course"
	quiz_route "clean-architecture/api/router/quiz"
	user_router "clean-architecture/api/router/user"
	"clean-architecture/bootstrap"
	"clean-architecture/infrastructor/mongo"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"time"
)

func SetUp(env *bootstrap.Database, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	publicRouter := gin.Group("")
	privateRouter := gin.Group("/admin")

	// Middleware
	publicRouter.Use(
		middleware.CORSPublic(),
		middleware.RateLimiter(),
		middleware.Recover(),
		gzip.Gzip(gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{",*"})),
	)

	privateRouter.Use(
		middleware.CORS(),
		middleware.RateLimiter(),
		middleware.Recover(),
		gzip.Gzip(gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{",*"})),
		//middleware.DeserializeUser(),
	)

	// This is a CORS method for check IP validation
	publicRouter.OPTIONS("/*path", middleware.OptionMessage)

	// All Public APIs
	user_router.GoogleAuthRouter(env, timeout, db, publicRouter)
	user_router.RefreshTokenRouter(env, timeout, db, publicRouter)
	quiz_route.QuizRouter(env, timeout, db, publicRouter)
	course_route.CourseRouter(env, timeout, db, publicRouter)
	publicRouter.GET("/logout", middleware.DeserializeUser(), user_controller.LogoutUser)

	// All Private API
	quiz_route.AdminQuizRouter(env, timeout, db, privateRouter)
	audio_route.AdminAudioRouter(env, timeout, db, privateRouter)
	course_route.AdminCourseRouter(env, timeout, db, privateRouter)
}
