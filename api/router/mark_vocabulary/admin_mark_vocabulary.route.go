package mark_vocabulary_route

import (
	mark_vocabulary_route "clean-architecture/api/controller/mark_vocabulary"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	mark_list_domain "clean-architecture/domain/mark_list"
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	admin_repository "clean-architecture/repository/admin"
	mark_vacabulary_repository "clean-architecture/repository/mark_vacabulary"
	admin_usecase "clean-architecture/usecase/admin"
	mark_vacabulary_usecase "clean-architecture/usecase/mark_vocabulary"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminMarkVocabularyRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ma := mark_vacabulary_repository.NewMarkVocabularyRepository(db, mark_list_domain.CollectionMarkList, vocabulary_domain.CollectionVocabulary, mark_vocabulary_domain.CollectionMark)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	markVocabulary := &mark_vocabulary_route.MarkVocabularyController{
		MarkVocabularyUseCase: mark_vacabulary_usecase.NewMarkVocabularyUseCase(ma, timeout),
		AdminUseCase:          admin_usecase.NewAdminUseCase(ad, timeout),
		Database:              env,
	}

	router := group.Group("/mark_vocabulary")
	router.GET("/fetch/mark_list_id", middleware.DeserializeUser(), markVocabulary.FetchManyByMarkListIdInAdmin)
}
