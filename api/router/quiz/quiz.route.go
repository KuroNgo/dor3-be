package quiz_route

import (
	quiz_controller "clean-architecture/api/controller/quiz"
	"clean-architecture/bootstrap"
	lesson_domain "clean-architecture/domain/lesson"
	quiz_domain "clean-architecture/domain/quiz"
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	quiz_repository "clean-architecture/repository/quiz"
	quiz_usecase "clean-architecture/usecase/quiz"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func QuizRouter(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	qu := quiz_repository.NewQuizRepository(db, quiz_domain.CollectionQuiz, lesson_domain.CollectionLesson, unit_domain.CollectionUnit, vocabulary_domain.CollectionVocabulary)
	quiz := &quiz_controller.QuizController{
		QuizUseCase: quiz_usecase.NewQuizUseCase(qu, timeout),
		Database:    env,
	}

	router := group.Group("/quiz")
	router.GET("/fetch", quiz.FetchManyQuiz)
}
