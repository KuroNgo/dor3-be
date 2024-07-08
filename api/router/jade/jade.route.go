package jade_route

import (
	jade_controller "clean-architecture/api/controller/jade"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	jade_domain "clean-architecture/domain/jade"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	admin_repository "clean-architecture/repository/admin"
	jade_repository "clean-architecture/repository/jade"
	user_repository "clean-architecture/repository/user"
	admin_usecase "clean-architecture/usecase/admin"
	jade_usecase "clean-architecture/usecase/jade"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func JadeRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	jd := jade_repository.NewJadeRepository(db, jade_domain.CollectionJade)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	jade := &jade_controller.JadeController{
		JadeUseCase:  jade_usecase.NewJadeUseCase(jd, timeout),
		AdminUseCase: admin_usecase.NewAdminUseCase(ad, timeout),
		UserUseCase:  usecase.NewUserUseCase(ur, timeout),
		Database:     env,
	}

	router := group.Group("/jade")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch", jade.FetchJadeInUser)
	router.GET("/fetch/rank", jade.RankInUser)
	router.POST("/create", jade.CreateJadeInUser)
	router.PUT("/update", jade.UpdateJade)
}
