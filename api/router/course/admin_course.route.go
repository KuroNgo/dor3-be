package course_route

import (
	course_controller "clean-architecture/api/controller/course"
	"clean-architecture/bootstrap"
	course_domain "clean-architecture/domain/course"
	"clean-architecture/infrastructor/mongo"
	course_repository "clean-architecture/repository/course"
	course_usecase "clean-architecture/usecase/course"
	"github.com/gin-gonic/gin"
	"time"
)

func AdminCourseRouter(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	co := course_repository.NewCourseRepository(db, course_domain.CollectionCourse)
	course := &course_controller.CourseController{
		CourseUseCase: course_usecase.NewCourseUseCase(co, timeout),
		Database:      env,
	}

	router := group.Group("/course")
	router.POST("/create", course.CreateOneCourse)

}
