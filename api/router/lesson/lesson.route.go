package lesson_route

import (
	lesson_controller "clean-architecture/api/controller/lesson"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	course_domain "clean-architecture/domain/course"
	image_domain "clean-architecture/domain/image"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	image_repository "clean-architecture/repository/image"
	lesson_repository "clean-architecture/repository/lesson"
	user_repository "clean-architecture/repository/user"
	image_usecase "clean-architecture/usecase/image"
	lesson_usecase "clean-architecture/usecase/lesson"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func LessonRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	le := lesson_repository.NewLessonRepository(db, lesson_domain.CollectionLesson, lesson_domain.CollectionLessonProcess, course_domain.CollectionCourse, unit_domain.CollectionUnit, vocabulary_domain.CollectionVocabulary)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)
	im := image_repository.NewImageRepository(db, image_domain.CollectionImage)

	lesson := &lesson_controller.LessonController{
		LessonUseCase: lesson_usecase.NewLessonUseCase(le, timeout),
		ImageUseCase:  image_usecase.NewImageUseCase(im, timeout),
		UserUseCase:   usecase.NewUserUseCase(ur, timeout),
		Database:      env,
	}

	router := group.Group("/lesson")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch", lesson.FetchManyInUser)
	router.GET("/fetch/course_id", lesson.FetchByIdCourseInUser)
	router.GET("/fetch/_id", lesson.FetchByIdInUser)
	router.GET("/fetch/not", lesson.FetchManyNotPaginationInUser)
}
