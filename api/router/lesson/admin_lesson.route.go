package lesson_route

import (
	lesson_controller "clean-architecture/api/controller/lesson"
	"clean-architecture/bootstrap"
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	"clean-architecture/infrastructor/mongo"
	lesson_repository "clean-architecture/repository/lesson"
	lesson_usecase "clean-architecture/usecase/lesson"
	"github.com/gin-gonic/gin"
	"time"
)

func AdminLessonRoute(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	le := lesson_repository.NewLessonRepository(db, lesson_domain.CollectionLesson, course_domain.CollectionCourse, unit_domain.CollectionUnit)
	lesson := &lesson_controller.LessonController{
		LessonUseCase: lesson_usecase.NewLessonUseCase(le, timeout),
		Database:      env,
	}

	router := group.Group("/lesson")
	router.POST("/create", lesson.CreateOneLesson)
	router.POST("/create/file", lesson.CreateLessonWithFile)
	router.POST("/upsert", lesson.UpsertOneLesson)
	router.PUT("/update", lesson.UpdateOneLesson)
	router.DELETE("/delete", lesson.DeleteOneLesson)
}
