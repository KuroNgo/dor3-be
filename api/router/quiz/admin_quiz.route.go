package quiz_route

import (
	quiz_controller "clean-architecture/api/controller/quiz"
	"clean-architecture/bootstrap"
	quiz_domain "clean-architecture/domain/quiz"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/infrastructor/mongo"
	quiz_repository "clean-architecture/repository/quiz"
	user_repository "clean-architecture/repository/user"
	quiz_usecase "clean-architecture/usecase/quiz"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"time"
)

func AdminQuizRouter(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	qu := quiz_repository.NewQuizRepository(db, quiz_domain.CollectionQuiz)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	quiz := &quiz_controller.QuizController{
		QuizUseCase: quiz_usecase.NewQuizUseCase(qu, timeout),
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
		Database:    env,
	}

	router := group.Group("/quiz")
	router.POST("/create", quiz.CreateOneQuiz)
	router.PUT("/update/:_id", quiz.UpdateOneQuiz)
	router.DELETE("/delete/:_id", quiz.DeleteOneQuiz)
}
