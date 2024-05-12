package mark_list_route

import (
	marklist_controller "clean-architecture/api/controller/mark_list"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	mark_list_domain "clean-architecture/domain/mark_list"
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	user_domain "clean-architecture/domain/user"
	mark_list_repository "clean-architecture/repository/mark_list"
	user_repository "clean-architecture/repository/user"
	mark_list_usecase "clean-architecture/usecase/mark_list"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func MarkListRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ma := mark_list_repository.NewListRepository(db, mark_list_domain.CollectionMarkList, mark_vocabulary_domain.CollectionMark)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	markList := &marklist_controller.MarkListController{
		MarkListUseCase: mark_list_usecase.NewMarkListUseCase(ma, timeout),
		UserUseCase:     usecase.NewUserUseCase(ur, timeout),
		Database:        env,
	}

	router := group.Group("/mark_list")
	router.GET("/fetch", middleware.DeserializeUser(), markList.FetchManyMarkListByUserID)
	router.POST("/create", middleware.DeserializeUser(), markList.CreateOneMarkList)
	router.PATCH("/update", middleware.DeserializeUser(), markList.UpdateOneMarkList)
	router.DELETE("/delete/_id", middleware.DeserializeUser(), markList.DeleteOneMarkList)
}
