package exam_options_route

import (
	exam_options_controller "clean-architecture/api/controller/exam_options"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exam_options_domain "clean-architecture/domain/exam_options"
	exam_question_domain "clean-architecture/domain/exam_question"
	user_domain "clean-architecture/domain/user"
	admin_repository "clean-architecture/repository/admin"
	exam_options_repository "clean-architecture/repository/exam_options"
	admin_usecase "clean-architecture/usecase/admin"
	exam_options_usecase "clean-architecture/usecase/exam_options"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminExamOptionsRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	opt := exam_options_repository.NewExamOptionsRepository(db, exam_question_domain.CollectionExamQuestion, exam_options_domain.CollectionExamOptions)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	options := &exam_options_controller.ExamOptionsController{
		ExamOptionsUseCase: exam_options_usecase.NewExamOptionsUseCase(opt, timeout),
		AdminUseCase:       admin_usecase.NewAdminUseCase(ad, timeout),
		Database:           env,
	}

	router := group.Group("/exam/options")
	router.POST("/create", options.CreateOneExamOptions)
	router.PATCH("/update", options.UpdateOneExamOptions)
	router.DELETE("/delete/:_id", options.DeleteOneExamOptions)
}
