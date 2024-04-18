package exam_route

import (
	exam_controller "clean-architecture/api/controller/exam"
	"clean-architecture/bootstrap"
	exam_domain "clean-architecture/domain/exam"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	exam_repository "clean-architecture/repository/exam"
	user_repository "clean-architecture/repository/user"
	exam_usecase "clean-architecture/usecase/exam"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExamRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ex := exam_repository.NewExamRepository(db, exam_domain.CollectionExam, lesson_domain.CollectionLesson, unit_domain.CollectionUnit, vocabulary_domain.CollectionVocabulary)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	exam := &exam_controller.ExamsController{
		ExamUseCase: exam_usecase.NewExamUseCase(ex, timeout),
		UserUseCase: user_usecase.NewUserUseCase(users, timeout),
		Database:    env,
	}

	router := group.Group("/exam")
	router.GET("/fetch", exam.FetchManyExam)
}
