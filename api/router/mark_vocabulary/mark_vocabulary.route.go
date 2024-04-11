package mark_vocabulary_route

import (
	mark_vocabulary_route "clean-architecture/api/controller/mark_vocabulary"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	mark_list_domain "clean-architecture/domain/mark_list"
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	"clean-architecture/infrastructor/mongo"
	mark_vacabulary_repository "clean-architecture/repository/mark_vacabulary"
	user_repository "clean-architecture/repository/user"
	mark_vacabulary_usecase "clean-architecture/usecase/mark_vocabulary"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"time"
)

func MarkVocabularyRoute(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ma := mark_vacabulary_repository.NewMarkVocabularyRepository(db, mark_list_domain.CollectionMarkList, vocabulary_domain.CollectionVocabulary, mark_vocabulary_domain.CollectionMark)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	markVocabulary := &mark_vocabulary_route.MarkVocabularyController{
		MarkVocabularyUseCase: mark_vacabulary_usecase.NewMarkVocabularyUseCase(ma, timeout),
		UserUseCase:           usecase.NewUserUseCase(ur, timeout),
		Database:              env,
	}

	router := group.Group("/mark_list")
	router.GET("/fetch/:_id", middleware.DeserializeUser(), markVocabulary.FetchManyByMarkListIdAndUserId)
	router.POST("/create/:mark_list_id/:vocabulary_id", middleware.DeserializeUser(), markVocabulary.CreateOneMarkVocabulary)
	router.DELETE("/delete/:_id", middleware.DeserializeUser(), markVocabulary.DeleteOneMarkVocabulary)
}
