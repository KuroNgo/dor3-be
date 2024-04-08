package exercise_route

import (
	exercise_controller "clean-architecture/api/controller/exercise"
	"clean-architecture/bootstrap"
	exercise_domain "clean-architecture/domain/exercise"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/infrastructor/mongo"
	exercise_repository "clean-architecture/repository/exercise"
	user_repository "clean-architecture/repository/user"
	exercise_usecase "clean-architecture/usecase/exercise"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"time"
)

func ExerciseRoute(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ex := exercise_repository.NewExerciseRepository(db, vocabulary_domain.CollectionVocabulary, exercise_domain.CollectionExercise)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	exercise := &exercise_controller.ExerciseController{
		ExerciseUseCase: exercise_usecase.NewExerciseUseCase(ex, timeout),
		UserUseCase:     usecase.NewUserUseCase(ur, timeout),
		Database:        env,
	}

	router := group.Group("/exercise")
	router.GET("/fetch", exercise.FetchMany)
}
