package user_router

import (
	user_controller "clean-architecture/api/controller/user"
	"clean-architecture/bootstrap"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/infrastructor/mongo"
	user_repository "clean-architecture/repository/user"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"time"
)

func RefreshTokenRouter(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)
	token := &user_controller.RefreshTokenController{
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
		Database:    env,
	}

	router := group.Group("/token")
	router.GET("/refresh", token.RefreshToken)
}
