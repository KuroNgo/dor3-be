package exercise_route

import (
	exercise_controller "clean-architecture/api/controller/exercise"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exercise_domain "clean-architecture/domain/exercise"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	admin_repository "clean-architecture/repository/admin"
	exercise_repository "clean-architecture/repository/exercise"
	admin_usecase "clean-architecture/usecase/admin"
	exercise_usecase "clean-architecture/usecase/exercise"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminExerciseRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ex := exercise_repository.NewExerciseRepository(db, lesson_domain.CollectionLesson, unit_domain.CollectionUnit, vocabulary_domain.CollectionVocabulary, exercise_domain.CollectionExercise, exercise_domain.CollectionExercise)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	exercise := &exercise_controller.ExerciseController{
		ExerciseUseCase: exercise_usecase.NewExerciseUseCase(ex, timeout),
		AdminUseCase:    admin_usecase.NewAdminUseCase(ad, timeout),
		Database:        env,
	}

	router := group.Group("/exercise")
	router.POST("/create", exercise.CreateOneExercise)
	router.PATCH("/update", exercise.UpdateOneExercise)
	router.DELETE("/delete/:_id", exercise.DeleteOneExercise)
}
