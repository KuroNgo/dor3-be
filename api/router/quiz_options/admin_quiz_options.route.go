package quiz_options_route

import (
	quiz_options_controller "clean-architecture/api/controller/quiz_options"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exercise_domain "clean-architecture/domain/exercise"
	quiz_options_domain "clean-architecture/domain/quiz_options"
	user_domain "clean-architecture/domain/user"
	admin_repository "clean-architecture/repository/admin"
	quiz_options_repository "clean-architecture/repository/quiz_options"
	admin_usecase "clean-architecture/usecase/admin"
	"clean-architecture/usecase/quiz_options"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminQuizOptionsRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	opt := quiz_options_repository.NewQuizOptionsRepository(db, exercise_domain.CollectionExercise, quiz_options_domain.CollectionQuizOptions)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	options := &quiz_options_controller.QuizOptionsController{
		QuizOptionsUseCase: quiz_options_usecase.NewQuizOptionsUseCase(opt, timeout),
		AdminUseCase:       admin_usecase.NewAdminUseCase(ad, timeout),
		Database:           env,
	}

	router := group.Group("/quiz/options")
	router.POST("/create", options.CreateOneQuizOptions)
	router.PATCH("/update", options.UpdateOneQuizOptions)
	router.DELETE("/delete/:_id", options.DeleteOneQuizOptions)
}
