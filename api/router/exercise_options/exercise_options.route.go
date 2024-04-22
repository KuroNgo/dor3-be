package exercise_options_route

import (
	exercise_options_controller "clean-architecture/api/controller/exercise_options"
	"clean-architecture/bootstrap"
	exercise_options_domain "clean-architecture/domain/exercise_options"
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
	user_domain "clean-architecture/domain/user"
	exercise_options_repository "clean-architecture/repository/exercise_options"
	user_repository "clean-architecture/repository/user"
	exercise_options_usecase "clean-architecture/usecase/exercise_options"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExerciseOptionsRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	opt := exercise_options_repository.NewExamOptionsRepository(db, exercise_questions_domain.CollectionExerciseQuestion, exercise_options_domain.CollectionExerciseOptions)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	options := &exercise_options_controller.ExerciseOptionsController{
		ExerciseOptionsUseCase: exercise_options_usecase.NewExerciseOptionsUseCase(opt, timeout),
		UserUseCase:            usecase.NewUserUseCase(ur, timeout),
		Database:               env,
	}

	router := group.Group("/exercise/options")
	router.GET("/fetch", options.FetchManyExerciseOptions)
}
