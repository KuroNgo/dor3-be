package exercise_result_route

import (
	exercise_result_controller "clean-architecture/api/controller/exercise_result"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exercise_domain "clean-architecture/domain/exercise"
	exercise_result_domain "clean-architecture/domain/exercise_result"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/repository/exercise_result"
	user_repository "clean-architecture/repository/user"
	exercise_result_usecase "clean-architecture/usecase/exercise_result"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExerciseResultRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	res := exercise_result_repository.NewExerciseResultRepository(db, exercise_result_domain.CollectionExerciseResult, exercise_domain.CollectionExercise)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	result := &exercise_result_controller.ExerciseResultController{
		ExerciseResultUseCase: exercise_result_usecase.NewExerciseQuestionUseCase(res, timeout),
		UserUseCase:           user_usecase.NewUserUseCase(users, timeout),
		Database:              env,
	}

	router := group.Group("/exercise/result")
	router.GET("/fetch/exercise_id", middleware.DeserializeUser(), result.GetResultsByUserIDAndExerciseID)
	router.GET("/fetch/0/exercise_id", middleware.DeserializeUser(), result.FetchResultByExerciseID)
	router.POST("/create", middleware.DeserializeUser(), result.CreateOneExercise)
	router.DELETE("/delete/_id", middleware.DeserializeUser(), result.DeleteOneExercise)
}
