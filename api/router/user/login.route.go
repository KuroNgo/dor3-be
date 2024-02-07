package user_router

import (
	controller "clean-architecture/api/controller/user"
	"clean-architecture/bootstrap"
	user_domain "clean-architecture/domain/request/user"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/repository"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"time"
)

func NewLoginRouter(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, user_domain.CollectionUser)
	lc := &controller.LoginController{
		LoginUseCase: usecase.NewLoginUseCase(ur, timeout),
		Env:          env,
	}
	group.POST("/username", lc.LoginByUserName)
	group.POST("email", lc.LoginByEmail)
}
