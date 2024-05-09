package exam_answer_route

import (
	exam_answer_controller "clean-architecture/api/controller/exam_answer"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exam_domain "clean-architecture/domain/exam"
	exam_answer_domain "clean-architecture/domain/exam_answer"
	exam_question_domain "clean-architecture/domain/exam_question"
	user_domain "clean-architecture/domain/user"
	exam_answer_repository "clean-architecture/repository/exam_answer"
	user_repository "clean-architecture/repository/user"
	exam_answer_usecase "clean-architecture/usecase/exam_answer"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExamAnswerRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ans := exam_answer_repository.NewExamAnswerRepository(db, exam_question_domain.CollectionExamQuestion, exam_answer_domain.CollectionExamAnswers, exam_domain.CollectionExam)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	answer := &exam_answer_controller.ExamAnswerController{
		ExamAnswerUseCase: exam_answer_usecase.NewExamAnswerUseCase(ans, timeout),
		UserUseCase:       user_usecase.NewUserUseCase(users, timeout),
		Database:          env,
	}

	router := group.Group("/exam/answer")
	router.GET("/fetch", middleware.DeserializeUser(), answer.FetchManyAnswerByUserIDAndQuestionID)
	router.POST("/create", middleware.DeserializeUser(), answer.CreateOneExamAnswer)
	router.DELETE("/1/delete", middleware.DeserializeUser(), answer.DeleteOneAnswer)
	router.DELETE("/all/delete", middleware.DeserializeUser(), answer.DeleteAllAnswerInExamID)
}
