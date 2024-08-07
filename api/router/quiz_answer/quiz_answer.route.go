package quiz_answer_route

import (
	quiz_answer_controller "clean-architecture/api/controller/quiz_answer"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	quiz_domain "clean-architecture/domain/quiz"
	quiz_answer_domain "clean-architecture/domain/quiz_answer"
	quiz_question_domain "clean-architecture/domain/quiz_question"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	quiz_answer_repository "clean-architecture/repository/quiz_answer"
	user_repository "clean-architecture/repository/user"
	quiz_answer_usecase "clean-architecture/usecase/quiz_answer"

	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func QuizAnswerRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ans := quiz_answer_repository.NewQuizAnswerRepository(db, quiz_question_domain.CollectionQuizQuestion, quiz_answer_domain.CollectionQuizAnswers, quiz_domain.CollectionQuiz)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	answer := &quiz_answer_controller.QuizAnswerController{
		QuizAnswerUseCase: quiz_answer_usecase.NewQuizResultUseCase(ans, timeout),
		UserUseCase:       usecase.NewUserUseCase(ur, timeout),
		Database:          env,
	}

	router := group.Group("/quiz/answer")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch", answer.FetchManyAnswerByUserIDAndQuestionID)
	router.POST("/create", answer.CreateOneQuizAnswer)
	router.DELETE("/delete", answer.DeleteOneAnswer)
	router.DELETE("/all/delete", answer.DeleteAllAnswerInQuizID)
}
