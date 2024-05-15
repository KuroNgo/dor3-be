package quiz_route

import (
	quiz_controller "clean-architecture/api/controller/quiz"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	lesson_domain "clean-architecture/domain/lesson"
	quiz_domain "clean-architecture/domain/quiz"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	quiz_repository "clean-architecture/repository/quiz"
	user_repository "clean-architecture/repository/user"
	quiz_usecase "clean-architecture/usecase/quiz"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func QuizRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	qu := quiz_repository.NewQuizRepository(db, quiz_domain.CollectionQuiz, lesson_domain.CollectionLesson, unit_domain.CollectionUnit)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	quiz := &quiz_controller.QuizController{
		QuizUseCase: quiz_usecase.NewQuizUseCase(qu, timeout),
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
		Database:    env,
	}

	router := group.Group("/quiz")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch", quiz.FetchManyQuiz)
	router.GET("/n/fetch", quiz.FetchManyQuizByUnitID)
	router.GET("/1/fetch", quiz.FetchOneQuizByUnitID)
}
