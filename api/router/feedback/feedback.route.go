package feedback_route

import (
	feedback_controller "clean-architecture/api/controller/feedback"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	feedback_domain "clean-architecture/domain/feedback"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	feedback_repository "clean-architecture/repository/feedback"
	user_repository "clean-architecture/repository/user"
	feedback_usecase "clean-architecture/usecase/feedback"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func FeedbackRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	fe := feedback_repository.NewFeedbackRepository(db, feedback_domain.CollectionFeedback, user_domain.CollectionUser)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	feedback := &feedback_controller.FeedbackController{
		FeedbackUseCase: feedback_usecase.NewFeedbackUseCase(fe, timeout),
		UserUseCase:     usecase.NewUserUseCase(ur, timeout),
		Database:        env,
	}

	router := group.Group("/feedback")
	router.Use(middleware.DeserializeUser())
	router.POST("/create", feedback.CreateOneFeedback)
}
