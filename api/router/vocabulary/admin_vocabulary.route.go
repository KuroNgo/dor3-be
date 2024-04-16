package vocabulary_route

import (
	vocabulary_controller "clean-architecture/api/controller/vocabulary"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	mark_domain "clean-architecture/domain/mark_vocabulary"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	admin_repository "clean-architecture/repository/admin"
	vocabulary_repository "clean-architecture/repository/vocabulary"
	admin_usecase "clean-architecture/usecase/admin"
	vocabulary_usecase "clean-architecture/usecase/vocabulary"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminVocabularyRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	vo := vocabulary_repository.NewVocabularyRepository(db, vocabulary_domain.CollectionVocabulary, mark_domain.CollectionMark, unit_domain.CollectionUnit)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)
	vocabulary := &vocabulary_controller.VocabularyController{
		VocabularyUseCase: vocabulary_usecase.NewVocabularyUseCase(vo, timeout),
		AdminUseCase:      admin_usecase.NewAdminUseCase(ad, timeout),
		Database:          env,
	}

	router := group.Group("/vocabulary")
	router.POST("/create", vocabulary.CreateOneVocabulary)
	router.POST("/create/file", vocabulary.CreateVocabularyWithFileInAdmin)
	router.POST("create/audio", vocabulary.GenerateVoice)
	router.PUT("/update/:_id", vocabulary.UpdateOneVocabulary)
	router.DELETE("/delete/:_id", vocabulary.DeleteOneVocabulary)
}
