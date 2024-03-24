package lesson_route

import (
	lesson_controller "clean-architecture/api/controller/lesson"
	"clean-architecture/bootstrap"
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	"clean-architecture/infrastructor/mongo"
	lesson_repository "clean-architecture/repository/lesson"
	lesson_usecase "clean-architecture/usecase/lesson"
	"github.com/gin-gonic/gin"
	"time"
)

func LessonRoute(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	le := lesson_repository.NewLessonRepository(db, lesson_domain.CollectionLesson, course_domain.CollectionCourse)
	lesson := &lesson_controller.LessonController{
		LessonUseCase: lesson_usecase.NewLessonUseCase(le, timeout),
		Database:      env,
	}

	router := group.Group("/lesson")
	router.GET("/fetch", lesson.FetchMany)
	router.GET("/fetch/:course-id", lesson.FetchByIdCourse)
}
