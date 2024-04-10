package router

import (
	"clean-architecture/api/middleware"
	admin_route "clean-architecture/api/router/admin"
	audio_route "clean-architecture/api/router/audio"
	course_route "clean-architecture/api/router/course"
	exercise_route "clean-architecture/api/router/exercise"
	image_route "clean-architecture/api/router/image"
	lesson_route "clean-architecture/api/router/lesson"
	mark_list_route "clean-architecture/api/router/mark_list"
	mark_vocabulary_route "clean-architecture/api/router/mark_vocabulary"
	quiz_route "clean-architecture/api/router/quiz"
	unit_route "clean-architecture/api/router/unit"
	user_route "clean-architecture/api/router/user"
	vocabulary_route "clean-architecture/api/router/vocabulary"
	"clean-architecture/bootstrap"
	"clean-architecture/infrastructor/mongo"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

func SetUp(env *bootstrap.Database, timeout time.Duration, db mongo.Database, gin *gin.Engine) {

	publicRouter := gin.Group("/api/")
	privateRouter := gin.Group("/api/admin")

	// Middleware
	publicRouter.Use(
		middleware.CORSPublic(),
		//middleware.RateLimiter(),
		middleware.Recover(),
		gzip.Gzip(gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{",*"})),
		middleware.StructuredLogger(&log.Logger),
	)

	privateRouter.Use(
		middleware.CORS(),
		//middleware.RateLimiter(),
		middleware.Recover(),
		gzip.Gzip(gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{",*"})),
		middleware.DeserializeUser(),
		middleware.StructuredLogger(&log.Logger),
	)

	// This is a CORS method for check IP validation
	publicRouter.OPTIONS("/*path", middleware.OptionMessage)

	// All Public APIs
	user_route.GoogleAuthRoute(env, timeout, db, publicRouter)
	user_route.UserRouter(env, timeout, db, publicRouter)
	exercise_route.ExerciseRoute(env, timeout, db, publicRouter)
	user_route.LoginFromRoleRoute(env, timeout, db, publicRouter)
	audio_route.AudioRoute(env, timeout, db, publicRouter)
	image_route.ImageRoute(env, timeout, db, publicRouter)
	quiz_route.QuizRouter(env, timeout, db, publicRouter)
	course_route.CourseRoute(env, timeout, db, publicRouter)
	lesson_route.LessonRoute(env, timeout, db, publicRouter)
	unit_route.UnitRouter(env, timeout, db, publicRouter)
	vocabulary_route.VocabularyRoute(env, timeout, db, publicRouter)

	// All Private API
	mark_vocabulary_route.MarkVocabularyRoute(env, timeout, db, publicRouter)
	mark_list_route.MarkListRoute(env, timeout, db, publicRouter)
	quiz_route.AdminQuizRouter(env, timeout, db, privateRouter)
	admin_route.AdminRouter(env, timeout, db, privateRouter)
	audio_route.AdminAudioRoute(env, timeout, db, privateRouter)
	image_route.AdminImageRoute(env, timeout, db, privateRouter)
	course_route.AdminCourseRoute(env, timeout, db, privateRouter)
	lesson_route.AdminLessonRoute(env, timeout, db, privateRouter)
	unit_route.AdminUnitRouter(env, timeout, db, privateRouter)
	vocabulary_route.AdminVocabularyRoute(env, timeout, db, privateRouter)
	exercise_route.AdminExerciseRoute(env, timeout, db, privateRouter)
}
