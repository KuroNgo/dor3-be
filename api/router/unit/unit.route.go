package unit_route

import (
	unit_controller "clean-architecture/api/controller/_unit"
	"clean-architecture/bootstrap"
	unit_domain "clean-architecture/domain/_unit"
	lesson_domain "clean-architecture/domain/lesson"
	"clean-architecture/infrastructor/mongo"
	unit_repo "clean-architecture/repository/unit"
	unit_usecase "clean-architecture/usecase/unit"
	"github.com/gin-gonic/gin"
	"time"
)

func UnitRouter(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	un := unit_repo.NewUnitRepository(db, unit_domain.CollectionUnit, lesson_domain.CollectionLesson)
	unit := &unit_controller.UnitController{
		UnitUseCase: unit_usecase.NewUnitUseCase(un, timeout),
		Database:    env,
	}

	router := group.Group("/unit")
	router.GET("/fetch", unit.FetchMany)
}
