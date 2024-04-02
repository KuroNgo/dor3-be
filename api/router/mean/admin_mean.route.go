package mean

import (
	mean_controller "clean-architecture/api/controller/mean"
	"clean-architecture/bootstrap"
	mean_domain "clean-architecture/domain/mean"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/infrastructor/mongo"
	mean_repository "clean-architecture/repository/mean"
	mean_usecase "clean-architecture/usecase/mean"
	"github.com/gin-gonic/gin"
	"time"
)

func AdminMeanRoute(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	me := mean_repository.NewMeanRepository(db, mean_domain.CollectionMean, vocabulary_domain.CollectionVocabulary)
	mean := &mean_controller.MeanController{
		MeanUseCase: mean_usecase.NewMeanUseCase(me, timeout),
		Database:    env,
	}

	router := group.Group("/lesson")
	router.POST("/create/file", mean.CreateMeanWithFile)
}
