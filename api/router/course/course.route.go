package course_route

import (
	course_controller "clean-architecture/api/controller/course"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	lesson_management_domain "clean-architecture/domain/user_process/lesson_management"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	course_repository "clean-architecture/repository/course"
	user_repository "clean-architecture/repository/user"
	course_usecase "clean-architecture/usecase/course"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func CourseRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	co := course_repository.NewCourseRepository(db, course_domain.CollectionCourse, lesson_management_domain.CollectionCourseProcess, lesson_domain.CollectionLesson, unit_domain.CollectionUnit, vocabulary_domain.CollectionVocabulary)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)
	course := &course_controller.CourseController{
		CourseUseCase: course_usecase.NewCourseUseCase(co, timeout),
		UserUseCase:   user_usecase.NewUserUseCase(users, timeout),
		Database:      env,
	}

	router := group.Group("/course")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch", course.FetchCourse)
	router.GET("/fetch/_id", course.FetchCourseByID)
}
