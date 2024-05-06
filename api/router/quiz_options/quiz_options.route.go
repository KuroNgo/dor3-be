package quiz_options_route

import (
	quiz_options_controller "clean-architecture/api/controller/quiz_options"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exercise_domain "clean-architecture/domain/exercise"
	quiz_options_domain "clean-architecture/domain/quiz_options"
	user_domain "clean-architecture/domain/user"
	quiz_options_repository "clean-architecture/repository/quiz_options"
	user_repository "clean-architecture/repository/user"
	quiz_options_usecase "clean-architecture/usecase/quiz_options"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func QuizOptionsRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	opt := quiz_options_repository.NewQuizOptionsRepository(db, exercise_domain.CollectionExercise, quiz_options_domain.CollectionQuizOptions)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	options := &quiz_options_controller.QuizOptionsController{
		QuizOptionsUseCase: quiz_options_usecase.NewQuizOptionsUseCase(opt, timeout),
		UserUseCase:        usecase.NewUserUseCase(ur, timeout),
		Database:           env,
	}

	router := group.Group("/quiz/options")
	router.GET("/fetch", middleware.DeserializeUser(), options.FetchManyQuizOptions)
}
