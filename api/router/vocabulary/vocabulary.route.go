package vocabulary_route

import (
	vocabulary_controller "clean-architecture/api/controller/vocabulary"
	"clean-architecture/bootstrap"
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/infrastructor/mongo"
	vocabulary_repository "clean-architecture/repository/vocabulary"
	vocabulary_usecase "clean-architecture/usecase/vocabulary"
	"github.com/gin-gonic/gin"
	"time"
)

func VocabularyRoute(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	vo := vocabulary_repository.NewVocabularyRepository(db, vocabulary_domain.CollectionVocabulary, unit_domain.CollectionUnit)
	vocabulary := &vocabulary_controller.VocabularyController{
		VocabularyUseCase: vocabulary_usecase.NewVocabularyUseCase(vo, timeout),
		Database:          env,
	}

	router := group.Group("/vocabulary")
	router.GET("/fetch", vocabulary.FetchMany)
	router.GET("/fetch-by-word", vocabulary.FetchByWord)
	router.GET("/fetch-by-lesson", vocabulary.FetchByLesson)
}
