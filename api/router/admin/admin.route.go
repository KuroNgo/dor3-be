package admin_route

import (
	admin_controller "clean-architecture/api/controller/admin"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/infrastructor/mongo"
	admin_repository "clean-architecture/repository/admin"
	admin_usecase "clean-architecture/usecase/admin"
	"github.com/gin-gonic/gin"
	"time"
)

func AdminRouter(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)
	admin := &admin_controller.AdminController{
		AdminUseCase: admin_usecase.NewAdminUseCase(ad, timeout),
		Database:     env,
	}

	router := group.Group("")
	router.POST("/signup", admin.SignUp)
	router.GET("/get-me", middleware.DeserializeUser(), admin.GetMe)
	router.GET("/refresh", admin.RefreshToken)
	router.GET("/logout", middleware.DeserializeUser(), admin.Logout)
}
