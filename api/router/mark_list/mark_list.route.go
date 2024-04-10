package mark_list

import (
	marklist_controller "clean-architecture/api/controller/marklist"
	"clean-architecture/bootstrap"
	mark_list_domain "clean-architecture/domain/mark_list"
	"clean-architecture/infrastructor/mongo"
	mark_list_repository "clean-architecture/repository/mark_list"
	mark_list_usecase "clean-architecture/usecase/mark_list"
	"github.com/gin-gonic/gin"
	"time"
)

func MarkListRoute(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ma := mark_list_repository.NewExerciseRepository(db, mark_list_domain.CollectionMarkList)
	markList := &marklist_controller.MarkListController{
		MarkListUseCase: mark_list_usecase.NewMarkListUseCase(ma, timeout),
		Database:        env,
	}

	router := group.Group("/mark_list")
	router.GET("/fetch", markList.FetchManyMarkList)
	router.POST("/create", markList.CreateOneMarkList)
	router.DELETE("/delete/:_id", markList.DeleteOneMarkList)
	router.PUT("/update/:_id", markList.UpdateOneMarkList)

}
