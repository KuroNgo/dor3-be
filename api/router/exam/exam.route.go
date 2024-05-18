package exam_route

import (
	exam_controller "clean-architecture/api/controller/exam"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exam_domain "clean-architecture/domain/exam"
	exam_question_domain "clean-architecture/domain/exam_question"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	admin_repository "clean-architecture/repository/admin"
	exam_repository "clean-architecture/repository/exam"
	user_repository "clean-architecture/repository/user"
	admin_usecase "clean-architecture/usecase/admin"
	exam_usecase "clean-architecture/usecase/exam"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExamRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ex := exam_repository.NewExamRepository(db, exam_domain.CollectionExam, lesson_domain.CollectionLesson, unit_domain.CollectionUnit, exam_question_domain.CollectionExamQuestion, vocabulary_domain.CollectionVocabulary)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	exam := &exam_controller.ExamsController{
		ExamUseCase:  exam_usecase.NewExamUseCase(ex, timeout),
		UserUseCase:  user_usecase.NewUserUseCase(users, timeout),
		AdminUseCase: admin_usecase.NewAdminUseCase(ad, timeout),
		Database:     env,
	}

	router := group.Group("/exam")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch/_id", exam.FetchOneExamByID)
	router.GET("/fetch", exam.FetchManyExam)
	router.GET("fetch/1/unit_id", exam.FetchOneExamByUnitID)
	router.GET("fetch/n/unit_id", exam.FetchManyExamByUnitID)
}
