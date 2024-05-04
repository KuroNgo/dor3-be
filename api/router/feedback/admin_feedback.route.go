package feedback_route

import (
	feedback_controller "clean-architecture/api/controller/feedback"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	feedback_domain "clean-architecture/domain/feedback"
	user_domain "clean-architecture/domain/user"
	admin_repository "clean-architecture/repository/admin"
	feedback_repository "clean-architecture/repository/feedback"
	admin_usecase "clean-architecture/usecase/admin"
	feedback_usecase "clean-architecture/usecase/feedback"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminFeedbackRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	fe := feedback_repository.NewFeedbackRepository(db, feedback_domain.CollectionFeedback)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	feedback := &feedback_controller.FeedbackController{
		FeedbackUseCase: feedback_usecase.NewFeedbackUseCase(fe, timeout),
		AdminUseCase:    admin_usecase.NewAdminUseCase(ad, timeout),
		Database:        env,
	}

	router := group.Group("/feedback")
	router.GET("/fetch", feedback.FetchMany)
	router.DELETE("/delete", feedback.DeleteOneFeedback)
}
