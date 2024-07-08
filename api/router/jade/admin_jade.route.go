package jade_route

import (
	jade_controller "clean-architecture/api/controller/jade"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	jade_domain "clean-architecture/domain/jade"
	user_domain "clean-architecture/domain/user"
	admin_repository "clean-architecture/repository/admin"
	jade_repository "clean-architecture/repository/jade"
	admin_usecase "clean-architecture/usecase/admin"
	jade_usecase "clean-architecture/usecase/jade"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminJadeRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	jd := jade_repository.NewJadeRepository(db, jade_domain.CollectionJade)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	jade := &jade_controller.JadeController{
		JadeUseCase:  jade_usecase.NewJadeUseCase(jd, timeout),
		AdminUseCase: admin_usecase.NewAdminUseCase(ad, timeout),
		Database:     env,
	}

	router := group.Group("/jade")
	router.GET("/fetch/rank", jade.RankInAdmin)
}
