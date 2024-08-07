package exercise_route

import (
	exercise_controller "clean-architecture/api/controller/exercise"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exercise_domain "clean-architecture/domain/exercise"
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	exercise_repository "clean-architecture/repository/exercise"
	user_repository "clean-architecture/repository/user"
	exercise_usecase "clean-architecture/usecase/exercise"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExerciseRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ex := exercise_repository.NewExerciseRepository(db, lesson_domain.CollectionLesson, unit_domain.CollectionUnit, vocabulary_domain.CollectionVocabulary, exercise_domain.CollectionExercise, exercise_questions_domain.CollectionExerciseQuestion)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	exercise := &exercise_controller.ExerciseController{
		ExerciseUseCase: exercise_usecase.NewExerciseUseCase(ex, timeout),
		UserUseCase:     usecase.NewUserUseCase(ur, timeout),
		Database:        env,
	}

	router := group.Group("/exercise")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch/_id", exercise.FetchOneExerciseByID)
	router.GET("/fetch/n/unit_id", exercise.FetchManyByUnitID)
	router.GET("/fetch", exercise.FetchManyExercise)
	router.GET("/fetch/1/unit_id", exercise.FetchOneByUnitID)
}
