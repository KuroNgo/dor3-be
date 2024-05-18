package quiz_question_route

import (
	quiz_question_controller "clean-architecture/api/controller/quiz_question"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	quiz_domain "clean-architecture/domain/quiz"
	quiz_question_domain "clean-architecture/domain/quiz_question"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	quiz_question_repository "clean-architecture/repository/quiz_question"
	user_repository "clean-architecture/repository/user"
	quiz_question_usecase "clean-architecture/usecase/quiz_question"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func QuizQuestionRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	quest := quiz_question_repository.NewQuizQuestionRepository(db, quiz_question_domain.CollectionQuizQuestion, quiz_domain.CollectionQuiz)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	question := &quiz_question_controller.QuizQuestionsController{
		QuizQuestionUseCase: quiz_question_usecase.NewQuizQuestionUseCase(quest, timeout),
		UserUseCase:         usecase.NewUserUseCase(ur, timeout),
		Database:            env,
	}

	router := group.Group("/quiz/question")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch", question.FetchManyQuizQuestion)
	router.GET("/fetch/_id", question.FetchOneQuizQuestionByID)
	router.GET("/fetch/1/quiz_id", question.FetchOneQuizQuestionByQuizID)
	router.GET("/fetch/n/quiz_id", question.FetchManyQuizQuestionByQuizID)
}
