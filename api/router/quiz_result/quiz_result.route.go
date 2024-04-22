package quiz_result_route

import (
	quiz_result_controller "clean-architecture/api/controller/quiz_result"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exercise_domain "clean-architecture/domain/exercise"
	quiz_question_domain "clean-architecture/domain/quiz_question"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/repository/quiz_result"
	user_repository "clean-architecture/repository/user"
	quiz_result_usecase "clean-architecture/usecase/quiz_result"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func QuizResultRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	res := quiz_result_repository.NewQuizResultRepository(db, quiz_question_domain.CollectionQuizQuestion, exercise_domain.CollectionExercise)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	result := &quiz_result_controller.QuizResultController{
		QuizResultUseCase: quiz_result_usecase.NewQuizQuestionUseCase(res, timeout),
		UserUseCase:       user_usecase.NewUserUseCase(users, timeout),
		Database:          env,
	}

	router := group.Group("/quiz/result")
	router.GET("/fetch/quiz/:quiz_id", middleware.DeserializeUser(), result.FetchResultByQuizID)
	router.GET("/fetch/user/:user_id/quiz/:quiz_id", middleware.DeserializeUser(), result.GetResultsByUserIDAndQuizID)
	router.POST("/create", middleware.DeserializeUser(), result.CreateOneQuizResult)
	router.DELETE("/delete/:_id", middleware.DeserializeUser(), result.DeleteOneQuizResult)
}
