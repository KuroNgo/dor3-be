package mark_vocabulary_route

import (
	mark_vocabulary_route "clean-architecture/api/controller/mark_vocabulary"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	mark_list_domain "clean-architecture/domain/mark_list"
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	mark_vacabulary_repository "clean-architecture/repository/mark_vacabulary"
	user_repository "clean-architecture/repository/user"
	vocabulary_repository "clean-architecture/repository/vocabulary"
	mark_vacabulary_usecase "clean-architecture/usecase/mark_vocabulary"
	usecase "clean-architecture/usecase/user"
	vocabulary_usecase "clean-architecture/usecase/vocabulary"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func MarkVocabularyRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ma := mark_vacabulary_repository.NewMarkVocabularyRepository(db, mark_list_domain.CollectionMarkList, vocabulary_domain.CollectionVocabulary, mark_vocabulary_domain.CollectionMark)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)
	vo := vocabulary_repository.NewVocabularyRepository(db, vocabulary_domain.CollectionVocabulary, mark_vocabulary_domain.CollectionMark, unit_domain.CollectionUnit)

	markVocabulary := &mark_vocabulary_route.MarkVocabularyController{
		MarkVocabularyUseCase: mark_vacabulary_usecase.NewMarkVocabularyUseCase(ma, timeout),
		VocabularyUseCase:     vocabulary_usecase.NewVocabularyUseCase(vo, timeout),
		UserUseCase:           usecase.NewUserUseCase(ur, timeout),
		Database:              env,
	}

	router := group.Group("/mark_vocabulary")
	router.GET("/fetch/mark_list_id", middleware.DeserializeUser(), markVocabulary.FetchManyByMarkListIdAndUserId)
	router.POST("/create", middleware.DeserializeUser(), markVocabulary.CreateOneMarkVocabulary)
	router.DELETE("/delete/_id", middleware.DeserializeUser(), markVocabulary.DeleteOneMarkVocabulary)
}
