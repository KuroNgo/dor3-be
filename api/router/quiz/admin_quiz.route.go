package quiz_route

import (
	quiz_controller "clean-architecture/api/controller/quiz"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	lesson_domain "clean-architecture/domain/lesson"
	quiz_domain "clean-architecture/domain/quiz"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	admin_repository "clean-architecture/repository/admin"
	quiz_repository "clean-architecture/repository/quiz"
	admin_usecase "clean-architecture/usecase/admin"
	quiz_usecase "clean-architecture/usecase/quiz"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminQuizRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	qu := quiz_repository.NewQuizRepository(db, quiz_domain.CollectionQuiz, lesson_domain.CollectionLesson, unit_domain.CollectionUnit)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	quiz := &quiz_controller.QuizController{
		QuizUseCase:  quiz_usecase.NewQuizUseCase(qu, timeout),
		AdminUseCase: admin_usecase.NewAdminUseCase(ad, timeout),
		Database:     env,
	}

	router := group.Group("/quiz")
	router.GET("/fetch", quiz.FetchManyQuizInAdmin)
	router.GET("/1/fetch", quiz.FetchOneQuizByUnitIDInAdmin)
	router.GET("/n/fetch", quiz.FetchManyQuizByUnitIDInAdmin)
	router.POST("/create", quiz.CreateOneQuiz)
	router.PATCH("/update", quiz.UpdateOneQuiz)
	router.DELETE("/delete/_id", quiz.DeleteOneQuiz)
}
