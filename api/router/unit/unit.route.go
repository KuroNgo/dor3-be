package unit_route

import (
	unit_controller "clean-architecture/api/controller/unit"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	unit_repo "clean-architecture/repository/unit"
	user_repository "clean-architecture/repository/user"
	unit_usecase "clean-architecture/usecase/unit"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func UnitRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	un := unit_repo.NewUnitRepository(db, unit_domain.CollectionUnit, lesson_domain.CollectionLesson, vocabulary_domain.CollectionVocabulary)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	unit := &unit_controller.UnitController{
		UnitUseCase: unit_usecase.NewUnitUseCase(un, timeout),
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
		Database:    env,
	}

	router := group.Group("/unit")
	router.GET("/fetch", unit.FetchMany)
	router.GET("/fetch/not", unit.FetchManyNotPagination)
	router.PATCH("/update/complete", middleware.DeserializeUser(), unit.UpdateCompleteUnit)
	router.GET("/fetch/lesson_id", unit.FetchByIdLesson)
}
