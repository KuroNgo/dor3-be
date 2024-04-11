package user_route

import (
	user_controller "clean-architecture/api/controller/user"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/infrastructor/mongo"
	user_repository "clean-architecture/repository/user"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"time"
)

func UserRouter(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)
	user := &user_controller.UserController{
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
		Database:    env,
	}

	router := group.Group("/user")
	router.GET("/get/mail", user.GetMail)
	router.POST("/signup", user.SignUp)
	router.GET("/info", middleware.DeserializeUser(), user.GetMe)
	router.GET("/refresh", user.RefreshToken)
	router.GET("/logout", middleware.DeserializeUser(), user.LogoutUser)
}
