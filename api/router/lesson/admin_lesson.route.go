package lesson_route

import (
	lesson_controller "clean-architecture/api/controller/lesson"
	"clean-architecture/bootstrap"
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	lesson_repository "clean-architecture/repository/lesson"
	user_repository "clean-architecture/repository/user"
	lesson_usecase "clean-architecture/usecase/lesson"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminLessonRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	le := lesson_repository.NewLessonRepository(db, lesson_domain.CollectionLesson, course_domain.CollectionCourse, unit_domain.CollectionUnit)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	lesson := &lesson_controller.LessonController{
		LessonUseCase: lesson_usecase.NewLessonUseCase(le, timeout),
		UserUseCase:   usecase.NewUserUseCase(ur, timeout),
		Database:      env,
	}

	router := group.Group("/lesson")
	router.POST("/create", lesson.CreateOneLesson)
	router.POST("/create/file", lesson.CreateLessonWithFile)
	router.POST("/upsert/:_id", lesson.UpsertOneLesson)
	router.PUT("/update/:_id", lesson.UpdateOneLesson)
	router.DELETE("/delete/:_id", lesson.DeleteOneLesson)
}
