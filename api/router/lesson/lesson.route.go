package lesson_route

import (
	lesson_controller "clean-architecture/api/controller/lesson"
	"clean-architecture/bootstrap"
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	lesson_repository "clean-architecture/repository/lesson"
	lesson_usecase "clean-architecture/usecase/lesson"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func LessonRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	le := lesson_repository.NewLessonRepository(db, lesson_domain.CollectionLesson, course_domain.CollectionCourse, unit_domain.CollectionUnit, vocabulary_domain.CollectionVocabulary)
	lesson := &lesson_controller.LessonController{
		LessonUseCase: lesson_usecase.NewLessonUseCase(le, timeout),
		Database:      env,
	}

	router := group.Group("/lesson")
	router.GET("/fetch", lesson.FetchMany)
	router.GET("/fetch/:course_id", lesson.FetchByIdCourse)
}
