package quiz_route

import (
	quiz_controller "clean-architecture/api/controller/quiz"
	"clean-architecture/bootstrap"
	quiz_domain "clean-architecture/domain/quiz"
	"clean-architecture/infrastructor/mongo"
	quiz_repository "clean-architecture/repository/quiz"
	quiz_usecase "clean-architecture/usecase/quiz"
	"github.com/gin-gonic/gin"
	"time"
)

func QuizRouter(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	qu := quiz_repository.NewQuizRepository(db, quiz_domain.CollectionQuiz)
	quiz := &quiz_controller.QuizController{
		QuizUseCase: quiz_usecase.NewQuizUseCase(qu, timeout),
		Database:    env,
	}

	router := group.Group("/quiz")
	router.GET("/fetch", quiz.FetchManyQuiz)
}
