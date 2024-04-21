package exam_options_route

import (
	exam_options_controller "clean-architecture/api/controller/exam_options"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exam_options_domain "clean-architecture/domain/exam_options"
	exam_question_domain "clean-architecture/domain/exam_question"
	user_domain "clean-architecture/domain/user"
	exam_options_repository "clean-architecture/repository/exam_options"
	user_repository "clean-architecture/repository/user"
	exam_options_usecase "clean-architecture/usecase/exam_options"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExamOptionsRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	opt := exam_options_repository.NewExamOptionsRepository(db, exam_question_domain.CollectionExamQuestion, exam_options_domain.CollectionExamOptions)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	options := &exam_options_controller.ExamOptionsController{
		ExamOptionsUseCase: exam_options_usecase.NewExamOptionsUseCase(opt, timeout),
		UserUseCase:        user_usecase.NewUserUseCase(users, timeout),
		Database:           env,
	}

	router := group.Group("/exam/options")
	router.GET("/fetch", middleware.DeserializeUser(), options.FetchManyExamOptions)
}
