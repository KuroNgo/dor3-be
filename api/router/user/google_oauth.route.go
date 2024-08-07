package user_route

import (
	user_controller "clean-architecture/api/controller/user"
	"clean-architecture/bootstrap"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	user_repository "clean-architecture/repository/user"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func GoogleAuthRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)
	ga := &user_controller.GoogleAuthController{
		GoogleAuthUseCase: usecase.NewGoogleUseCase(ur, timeout),
		Database:          env,
	}

	router := group.Group("/auth")
	router.GET("/google/callback", ga.GoogleLoginWithUser)

}
