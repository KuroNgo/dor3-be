package quiz_result_route

import (
	quiz_result_controller "clean-architecture/api/controller/quiz_result"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	quiz_domain "clean-architecture/domain/quiz"
	quiz_result_domain "clean-architecture/domain/quiz_result"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	"clean-architecture/repository/quiz_result"
	user_repository "clean-architecture/repository/user"
	quiz_result_usecase "clean-architecture/usecase/quiz_result"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func QuizResultRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	res := quiz_result_repository.NewQuizResultRepository(db, quiz_result_domain.CollectionQuizResult, quiz_domain.CollectionQuiz)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	result := &quiz_result_controller.QuizResultController{
		QuizResultUseCase: quiz_result_usecase.NewQuizQuestionUseCase(res, timeout),
		UserUseCase:       user_usecase.NewUserUseCase(users, timeout),
		Database:          env,
	}

	router := group.Group("/quiz/result")
	router.GET("/fetch/quiz_id", middleware.DeserializeUser(), result.FetchResultByQuizID)
	router.GET("/fetch", middleware.DeserializeUser(), result.GetResultsByUserIDAndQuizID)
	router.POST("/create", middleware.DeserializeUser(), result.CreateOneQuizResult)
	router.DELETE("/delete/_id", middleware.DeserializeUser(), result.DeleteOneQuizResult)
}
