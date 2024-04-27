package course_route

import (
	course_controller "clean-architecture/api/controller/course"
	"clean-architecture/bootstrap"
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	course_repository "clean-architecture/repository/course"
	course_usecase "clean-architecture/usecase/course"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func CourseRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	co := course_repository.NewCourseRepository(db, course_domain.CollectionCourse, lesson_domain.CollectionLesson, unit_domain.CollectionUnit, vocabulary_domain.CollectionVocabulary)
	course := &course_controller.CourseController{
		CourseUseCase: course_usecase.NewCourseUseCase(co, timeout),
		Database:      env,
	}

	router := group.Group("/course")
	router.GET("/fetch", course.FetchCourse)
	router.GET("/statistic", course.StatisticCourse)
}
