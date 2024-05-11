package vocabulary_route

import (
	vocabulary_controller "clean-architecture/api/controller/vocabulary"
	"clean-architecture/bootstrap"
	image_domain "clean-architecture/domain/image"
	mark_domain "clean-architecture/domain/mark_vocabulary"
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	image_repository "clean-architecture/repository/image"
	vocabulary_repository "clean-architecture/repository/vocabulary"
	image_usecase "clean-architecture/usecase/image"
	vocabulary_usecase "clean-architecture/usecase/vocabulary"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func VocabularyRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	vo := vocabulary_repository.NewVocabularyRepository(db, vocabulary_domain.CollectionVocabulary,
		mark_domain.CollectionMark, unit_domain.CollectionUnit)
	im := image_repository.NewImageRepository(db, image_domain.CollectionImage)

	vocabulary := &vocabulary_controller.VocabularyController{
		VocabularyUseCase: vocabulary_usecase.NewVocabularyUseCase(vo, timeout),
		ImageUseCase:      image_usecase.NewImageUseCase(im, timeout),
		Database:          env,
	}

	router := group.Group("/vocabulary")
	router.GET("/fetch", vocabulary.FetchMany)
	router.GET("/fetch/latest", vocabulary.FetchAllVocabularyLatest)
	router.GET("/fetch-all", vocabulary.FetchAllVocabulary)
	router.GET("/fetch-by-word", vocabulary.FetchByWord)
	router.GET("/fetch-by-lesson", vocabulary.FetchByLesson)
	router.GET("/fetch/unit", vocabulary.FetchByIdUnit)
}
