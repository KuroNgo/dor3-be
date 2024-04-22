package exercise_options_route

import (
	exercise_options_controller "clean-architecture/api/controller/exercise_options"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exercise_options_domain "clean-architecture/domain/exercise_options"
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
	user_domain "clean-architecture/domain/user"
	admin_repository "clean-architecture/repository/admin"
	exercise_options_repository "clean-architecture/repository/exercise_options"
	admin_usecase "clean-architecture/usecase/admin"
	"clean-architecture/usecase/exercise_options"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminExerciseOptionsRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	opt := exercise_options_repository.NewExamOptionsRepository(db, exercise_questions_domain.CollectionExerciseQuestion, exercise_options_domain.CollectionExerciseOptions)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	options := &exercise_options_controller.ExerciseOptionsController{
		ExerciseOptionsUseCase: exercise_options_usecase.NewExerciseOptionsUseCase(opt, timeout),
		AdminUseCase:           admin_usecase.NewAdminUseCase(ad, timeout),
		Database:               env,
	}

	router := group.Group("/exercise/options")
	router.POST("/create", options.CreateOneExerciseOptions)
	router.PATCH("/update", options.UpdateOneExerciseOptions)
	router.DELETE("/delete/:_id", options.DeleteOneExerciseOptions)
}
