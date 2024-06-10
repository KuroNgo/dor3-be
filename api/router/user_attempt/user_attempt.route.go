package user_attempt_route

import (
	user_attempt_controller "clean-architecture/api/controller/user_attempt"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exam_domain "clean-architecture/domain/exam"
	exercise_domain "clean-architecture/domain/exercise"
	quiz_domain "clean-architecture/domain/quiz"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	user_attempt_domain "clean-architecture/domain/user_process"
	user_repository "clean-architecture/repository/user"
	user_attempt_repository "clean-architecture/repository/user_attempt"
	user_usecase "clean-architecture/usecase/user"
	user_attempt_usecase "clean-architecture/usecase/user_attempt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func UserAttemptRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	userAttempt := user_attempt_repository.NewUserAttemptRepository(db, user_attempt_domain.CollectionUserExamManagement, exam_domain.CollectionExam,
		quiz_domain.CollectionQuiz, exercise_domain.CollectionExercise)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	result := &user_attempt_controller.UserAttemptController{
		UserAttemptUseCase: user_attempt_usecase.NewAttemptUseCase(userAttempt, timeout),
		UserUseCase:        user_usecase.NewUserUseCase(users, timeout),
		Database:           env,
	}

	router := group.Group("/user/attempt")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch", result.FetchManyResult)
}
