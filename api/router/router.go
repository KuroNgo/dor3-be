package router

import (
	"clean-architecture/api/middleware"
	activity_log_route "clean-architecture/api/router/activity_log"
	admin_route "clean-architecture/api/router/admin"
	course_route "clean-architecture/api/router/course"
	exam_route "clean-architecture/api/router/exam"
	exam_answer_route "clean-architecture/api/router/exam_answer"
	exam_question_route "clean-architecture/api/router/exam_question"
	exam_result_route "clean-architecture/api/router/exam_result"
	exercise_route "clean-architecture/api/router/exercise"
	exercise_answer_route "clean-architecture/api/router/exercise_answer"
	exercise_question_route "clean-architecture/api/router/exercise_question"
	exercise_result_route "clean-architecture/api/router/exercise_result"
	feedback_route "clean-architecture/api/router/feedback"
	image_route "clean-architecture/api/router/image"
	lesson_route "clean-architecture/api/router/lesson"
	mark_list_route "clean-architecture/api/router/mark_list"
	mark_vocabulary_route "clean-architecture/api/router/mark_vocabulary"
	quiz_route "clean-architecture/api/router/quiz"
	quiz_answer_route "clean-architecture/api/router/quiz_answer"
	quiz_question_route "clean-architecture/api/router/quiz_question"
	quiz_result_route "clean-architecture/api/router/quiz_result"
	unit_route "clean-architecture/api/router/unit"
	user_route "clean-architecture/api/router/user"
	user_attempt_route "clean-architecture/api/router/user_attempt"
	vocabulary_route "clean-architecture/api/router/vocabulary"
	"clean-architecture/bootstrap"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func SetUp(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, gin *gin.Engine) {
	value := activity_log_route.ActivityRoute(env, timeout, db)

	publicRouter := gin.Group("/api")
	privateRouter := gin.Group("/api/admin/")

	// Middleware
	publicRouter.Use(
		middleware.CORSPublic(),
		//middleware.RateLimiter(),
		middleware.Recover(),
		gzip.Gzip(gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{",*"})),
		middleware.StructuredLogger(&log.Logger, value),
	)

	privateRouter.Use(
		middleware.CORSPrivate(),
		//middleware.RateLimiter(),
		middleware.Recover(),
		gzip.Gzip(gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{",*"})),
		middleware.DeserializeUser(),
		middleware.StructuredLogger(&log.Logger, value),
	)

	// This is a CORS method for check IP validation
	publicRouter.OPTIONS("/*path", middleware.OptionMessage)

	// All Public APIs
	user_route.GoogleAuthRoute(env, timeout, db, publicRouter)
	user_route.UserRouter(env, timeout, db, publicRouter)
	admin_route.AdminRouter(env, timeout, db, publicRouter)
	user_route.LoginFromRoleRoute(env, timeout, db, publicRouter)
	user_attempt_route.UserAttemptRoute(env, timeout, db, publicRouter)

	feedback_route.FeedbackRoute(env, timeout, db, publicRouter)

	exam_route.ExamRoute(env, timeout, db, publicRouter)
	exam_answer_route.ExamAnswerRoute(env, timeout, db, publicRouter)
	exam_question_route.ExamQuestionRoute(env, timeout, db, publicRouter)
	exam_result_route.ExamResultRoute(env, timeout, db, publicRouter)

	exercise_route.ExerciseRoute(env, timeout, db, publicRouter)
	exercise_answer_route.ExerciseRoute(env, timeout, db, publicRouter)
	exercise_question_route.ExerciseQuestionRoute(env, timeout, db, publicRouter)
	exercise_result_route.ExerciseResultRoute(env, timeout, db, publicRouter)

	quiz_route.QuizRouter(env, timeout, db, publicRouter)
	quiz_answer_route.QuizAnswerRoute(env, timeout, db, publicRouter)
	quiz_question_route.QuizQuestionRoute(env, timeout, db, publicRouter)
	quiz_result_route.QuizResultRoute(env, timeout, db, publicRouter)

	mark_vocabulary_route.MarkVocabularyRoute(env, timeout, db, publicRouter)
	mark_list_route.MarkListRoute(env, timeout, db, publicRouter)

	image_route.ImageRoute(env, timeout, db, publicRouter)
	course_route.CourseRoute(env, timeout, db, publicRouter)
	lesson_route.LessonRoute(env, timeout, db, publicRouter)
	unit_route.UnitRouter(env, timeout, db, publicRouter)
	vocabulary_route.VocabularyRoute(env, timeout, db, publicRouter)

	// All Private API
	activity_log_route.AdminActivityRoute(env, timeout, db, privateRouter)
	feedback_route.AdminFeedbackRoute(env, timeout, db, privateRouter)

	exam_route.AdminExamRoute(env, timeout, db, privateRouter)
	exam_question_route.AdminExamQuestionRoute(env, timeout, db, privateRouter)

	exercise_route.AdminExerciseRoute(env, timeout, db, privateRouter)
	exercise_question_route.AdminExerciseQuestionRoute(env, timeout, db, privateRouter)

	quiz_route.AdminQuizRouter(env, timeout, db, privateRouter)
	quiz_question_route.AdminQuizQuestionRoute(env, timeout, db, privateRouter)

	image_route.AdminImageRoute(env, timeout, db, privateRouter)
	course_route.AdminCourseRoute(env, timeout, db, privateRouter)
	lesson_route.AdminLessonRoute(env, timeout, db, privateRouter)
	unit_route.AdminUnitRouter(env, timeout, db, privateRouter)
	vocabulary_route.AdminVocabularyRoute(env, timeout, db, privateRouter)

	mark_list_route.AdminMarkListRoute(env, timeout, db, privateRouter)
	mark_vocabulary_route.AdminMarkVocabularyRoute(env, timeout, db, privateRouter)
}
