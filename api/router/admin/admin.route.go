package admin_route

import (
	admin_controller "clean-architecture/api/controller/admin"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	admin_repository "clean-architecture/repository/admin"
	user_repository "clean-architecture/repository/user"
	admin_usecase "clean-architecture/usecase/admin"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	admin := &admin_controller.AdminController{
		AdminUseCase: admin_usecase.NewAdminUseCase(ad, timeout),
		UserUseCase:  usecase.NewUserUseCase(ur, timeout),
		Database:     env,
	}

	router := group.Group("/admin")
	router.POST("/signup", admin.SignUp)
	router.GET("/get-me", admin.GetMe)
	router.PUT("/update", admin.UpdateAdmin)
	router.GET("/refresh", admin.RefreshToken)
	router.GET("/logout", middleware.DeserializeUser(), admin.Logout)
	router.GET("/user/fetch", middleware.DeserializeUser(), admin.FetchManyUser)
	router.GET("/user/fetch/user_id", middleware.DeserializeUser(), admin.FetchUserByID)
}
