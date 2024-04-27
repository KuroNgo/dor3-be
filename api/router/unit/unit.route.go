package unit_route

import (
	unit_controller "clean-architecture/api/controller/unit"
	"clean-architecture/bootstrap"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	unit_repo "clean-architecture/repository/unit"
	unit_usecase "clean-architecture/usecase/unit"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func UnitRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	un := unit_repo.NewUnitRepository(db, unit_domain.CollectionUnit, lesson_domain.CollectionLesson, vocabulary_domain.CollectionVocabulary)

	unit := &unit_controller.UnitController{
		UnitUseCase: unit_usecase.NewUnitUseCase(un, timeout),
		Database:    env,
	}

	router := group.Group("/unit")
	router.GET("/fetch", unit.FetchMany)
	router.GET("/fetch/:lesson", unit.FetchByIdLesson)
}
