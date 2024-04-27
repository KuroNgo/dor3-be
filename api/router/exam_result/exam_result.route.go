package exam_result_route

import (
	exam_result_controller "clean-architecture/api/controller/exam_result"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exam_domain "clean-architecture/domain/exam"
	exam_result_domain "clean-architecture/domain/exam_result"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/repository/exam_result"
	user_repository "clean-architecture/repository/user"
	exam_result_usecase "clean-architecture/usecase/exam_result"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExamResultRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	res := exam_result_repository.NewExamResultRepository(db, exam_result_domain.CollectionExamResult, exam_domain.CollectionExam)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	result := &exam_result_controller.ExamResultController{
		ExamResultUseCase: exam_result_usecase.NewExamResultUseCase(res, timeout),
		UserUseCase:       user_usecase.NewUserUseCase(users, timeout),
		Database:          env,
	}

	router := group.Group("/exam/result")
	router.GET("/fetch/user_id/exam_id", middleware.DeserializeUser(), result.GetResultsByUserIDAndExamID)
	router.GET("/fetch/exam_id", middleware.DeserializeUser(), result.FetchResultByExamID)
	router.POST("/create", middleware.DeserializeUser(), result.CreateOneExamResult)
	router.DELETE("/delete/:_id", middleware.DeserializeUser(), result.DeleteOneExamResult)
}
