package unit_route

import (
	unit_controller "clean-architecture/api/controller/unit"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exam_domain "clean-architecture/domain/exam"
	exercise_domain "clean-architecture/domain/exercise"
	lesson_domain "clean-architecture/domain/lesson"
	quiz_domain "clean-architecture/domain/quiz"
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
	un := unit_repo.NewUnitRepository(db, unit_domain.CollectionUnit, lesson_domain.CollectionLesson,
		vocabulary_domain.CollectionVocabulary, exam_domain.CollectionExam, exercise_domain.CollectionExercise, quiz_domain.CollectionQuiz)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	unit := &unit_controller.UnitController{
		UnitUseCase: unit_usecase.NewUnitUseCase(un, timeout),
		UserUseCase: usecase.NewUserUseCase(ur, timeout),
		Database:    env,
	}

	router := group.Group("/unit")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch", unit.FetchMany)
	router.GET("/fetch/_id", unit.FetchById)
	router.GET("/fetch/not", unit.FetchManyNotPagination)
	router.GET("/fetch/lesson_id", unit.FetchByIdLesson)
}
