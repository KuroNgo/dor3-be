package mark_list_route

import (
	marklist_controller "clean-architecture/api/controller/mark_list"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	mark_list_domain "clean-architecture/domain/mark_list"
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	user_domain "clean-architecture/domain/user"
	admin_repository "clean-architecture/repository/admin"
	mark_list_repository "clean-architecture/repository/mark_list"
	admin_usecase "clean-architecture/usecase/admin"
	mark_list_usecase "clean-architecture/usecase/mark_list"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminMarkListRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ma := mark_list_repository.NewListRepository(db, mark_list_domain.CollectionMarkList, mark_vocabulary_domain.CollectionMark)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	markList := &marklist_controller.MarkListController{
		MarkListUseCase: mark_list_usecase.NewMarkListUseCase(ma, timeout),
		AdminUseCase:    admin_usecase.NewAdminUseCase(ad, timeout),
		Database:        env,
	}

	router := group.Group("/mark_list")
	router.GET("/fetch", middleware.DeserializeUser(), markList.FetchManyInAdmin)
	router.GET("/fetch/_id", middleware.DeserializeUser(), markList.FetchManyMarkListByIDInAdmin)
}
