package unit_route

import (
	unit_controller "clean-architecture/api/controller/unit"
	"clean-architecture/bootstrap"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/infrastructor/mongo"
	unit_repo "clean-architecture/repository/unit"
	unit_usecase "clean-architecture/usecase/unit"
	"github.com/gin-gonic/gin"
	"time"
)

func AdminUnitRouter(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	un := unit_repo.NewUnitRepository(db, unit_domain.CollectionUnit, lesson_domain.CollectionLesson, vocabulary_domain.CollectionVocabulary)
	unit := &unit_controller.UnitController{
		UnitUseCase: unit_usecase.NewUnitUseCase(un, timeout),
		Database:    env,
	}

	router := group.Group("/unit")
	router.POST("/create", unit.CreateOneUnit)
	router.POST("/create/file", unit.CreateUnitWithFile)
	router.PUT("/update/:_id", unit.UpdateOneUnit)
	router.POST("/upsert/:_id", unit.UpsertOneUnit)
	router.DELETE("/delete/:_id", unit.DeleteOneUnit)
}
