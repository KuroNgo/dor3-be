package user_route

import (
	user_controller "clean-architecture/api/controller/user"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/infrastructor/mongo"
	admin_repository "clean-architecture/repository/admin"
	user_repository "clean-architecture/repository/user"
	admin_usecase "clean-architecture/usecase/admin"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"time"
)

func LoginFromRoleRoute(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	login := &user_controller.LoginFromRoleController{
		UserUseCase:  usecase.NewUserUseCase(ur, timeout),
		AdminUseCase: admin_usecase.NewAdminUseCase(ad, timeout),
		Database:     env,
	}

	router := group.Group("/login")
	router.POST("/role", login.Login2)
}
