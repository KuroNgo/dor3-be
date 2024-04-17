package unit_route

import (
	unit_controller "clean-architecture/api/controller/unit"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	admin_repository "clean-architecture/repository/admin"
	unit_repo "clean-architecture/repository/unit"
	admin_usecase "clean-architecture/usecase/admin"
	unit_usecase "clean-architecture/usecase/unit"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminUnitRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	un := unit_repo.NewUnitRepository(db, unit_domain.CollectionUnit, lesson_domain.CollectionLesson, vocabulary_domain.CollectionVocabulary)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	unit := &unit_controller.UnitController{
		UnitUseCase:  unit_usecase.NewUnitUseCase(un, timeout),
		AdminUseCase: admin_usecase.NewAdminUseCase(ad, timeout),
		Database:     env,
	}

	router := group.Group("/unit")
	router.POST("/create", unit.CreateOneUnit)
	router.POST("/create/file", unit.CreateUnitWithFile)
	router.DELETE("/delete/:_id", unit.DeleteOneUnit)
}
