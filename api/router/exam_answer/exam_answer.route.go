package exam_answer_route

import (
	exam_answer_controller "clean-architecture/api/controller/exam_answer"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exam_domain "clean-architecture/domain/exam"
	exam_answer_domain "clean-architecture/domain/exam_answer"
	exam_options_domain "clean-architecture/domain/exam_options"
	exam_question_domain "clean-architecture/domain/exam_question"
	exam_result_domain "clean-architecture/domain/exam_result"
	exercise_domain "clean-architecture/domain/exercise"
	quiz_domain "clean-architecture/domain/quiz"
	user_domain "clean-architecture/domain/user"
	user_attempt_domain "clean-architecture/domain/user_attempt"
	user_detail_domain "clean-architecture/domain/user_detail"
	exam_answer_repository "clean-architecture/repository/exam_answer"
	exam_result_repository "clean-architecture/repository/exam_result"
	user_repository "clean-architecture/repository/user"
	user_attempt_repository "clean-architecture/repository/user_attempt"
	exam_answer_usecase "clean-architecture/usecase/exam_answer"
	exam_result_usecase "clean-architecture/usecase/exam_result"
	user_usecase "clean-architecture/usecase/user"
	user_attempt_usecase "clean-architecture/usecase/user_attempt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExamAnswerRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ans := exam_answer_repository.NewExamAnswerRepository(db, exam_question_domain.CollectionExamQuestion, exam_options_domain.CollectionExamOptions, exam_answer_domain.CollectionExamAnswers, exam_domain.CollectionExam)
	res := exam_result_repository.NewExamResultRepository(db, exam_result_domain.CollectionExamResult, exam_domain.CollectionExam)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)
	userAttempt := user_attempt_repository.NewUserAttemptRepository(db, user_attempt_domain.CollectionUserAttempt, exam_domain.CollectionExam, quiz_domain.CollectionQuiz, exercise_domain.CollectionExercise)

	answer := &exam_answer_controller.ExamAnswerController{
		ExamAnswerUseCase:  exam_answer_usecase.NewExamAnswerUseCase(ans, timeout),
		ExamResultUseCase:  exam_result_usecase.NewExamResultUseCase(res, timeout),
		UserUseCase:        user_usecase.NewUserUseCase(users, timeout),
		UserAttemptUseCase: user_attempt_usecase.NewAttemptUseCase(userAttempt, timeout),
		Database:           env,
	}

	router := group.Group("/exam/answer")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch", answer.FetchManyAnswerByUserIDAndQuestionID)
	router.POST("/create", answer.CreateOneExamAnswer)
	router.DELETE("/1/delete", answer.DeleteOneAnswer)
	router.DELETE("/all/delete", answer.DeleteAllAnswerInExamID)
}
