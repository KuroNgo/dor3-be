package user_route

import (
	user_controller "clean-architecture/api/controller/user"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	user_repository "clean-architecture/repository/user"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func UserRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)
	user := &user_controller.UserController{
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
		Database:    env,
	}

	router := group.Group("/user")
	router.POST("/signup", user.SignUp)
	router.PATCH("/update", middleware.DeserializeUser(), user.UpdateUser)
	router.PATCH("/verify", user.VerificationCode)
	router.PATCH("/verify/password", user.VerificationCodeForChangePassword)
	router.PATCH("/password/forget", user.ChangePassword)
	router.POST("/forget", user.ForgetPasswordInUser)
	router.GET("/info", middleware.DeserializeUser(), user.GetMe)
	router.GET("/refresh", middleware.DeserializeUser(), user.RefreshToken)
	router.GET("/logout", middleware.DeserializeUser(), user.LogoutUser)
}
