package exam_result_route

import (
	exam_result_controller "clean-architecture/api/controller/exam_result"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exam_domain "clean-architecture/domain/exam"
	exam_result_domain "clean-architecture/domain/exam_result"
	exercise_result_domain "clean-architecture/domain/exercise_result"
	quiz_domain "clean-architecture/domain/quiz"
	user_domain "clean-architecture/domain/user"
	user_attempt_domain "clean-architecture/domain/user_attempt"
	"clean-architecture/repository/exam_result"
	user_repository "clean-architecture/repository/user"
	user_attempt_repository "clean-architecture/repository/user_attempt"
	exam_result_usecase "clean-architecture/usecase/exam_result"
	user_usecase "clean-architecture/usecase/user"
	user_attempt_usecase "clean-architecture/usecase/user_attempt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExamResultRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	res := exam_result_repository.NewExamResultRepository(db, exam_result_domain.CollectionExamResult, exam_domain.CollectionExam)
	userAttempt := user_attempt_repository.NewUserAttemptRepository(db, user_attempt_domain.CollectionUserAttempt, exam_domain.CollectionExam, quiz_domain.CollectionQuiz, exercise_result_domain.CollectionExerciseResult)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	result := &exam_result_controller.ExamResultController{
		ExamResultUseCase:  exam_result_usecase.NewExamResultUseCase(res, timeout),
		UserAttemptUseCase: user_attempt_usecase.NewAttemptUseCase(userAttempt, timeout),
		UserUseCase:        user_usecase.NewUserUseCase(users, timeout),
		Database:           env,
	}

	router := group.Group("/exam/result")
	router.GET("/fetch/user_id/exam_id", middleware.DeserializeUser(), result.GetResultsByUserIDAndExamID)
	router.GET("/fetch/exam_id", middleware.DeserializeUser(), result.FetchResultByExamID)
	router.POST("/create", middleware.DeserializeUser(), result.CreateOneExamResult)
	router.DELETE("/delete/:_id", middleware.DeserializeUser(), result.DeleteOneExamResult)
}
