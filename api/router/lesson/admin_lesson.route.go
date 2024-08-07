package lesson_route

import (
	lesson_controller "clean-architecture/api/controller/lesson"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	course_domain "clean-architecture/domain/course"
	image_domain "clean-architecture/domain/image"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	admin_repository "clean-architecture/repository/admin"
	image_repository "clean-architecture/repository/image"
	lesson_repository "clean-architecture/repository/lesson"
	admin_usecase "clean-architecture/usecase/admin"
	image_usecase "clean-architecture/usecase/image"
	lesson_usecase "clean-architecture/usecase/lesson"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminLessonRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	le := lesson_repository.NewLessonRepository(db, lesson_domain.CollectionLesson, lesson_domain.CollectionLessonProcess, course_domain.CollectionCourse, unit_domain.CollectionUnit, vocabulary_domain.CollectionVocabulary)
	im := image_repository.NewImageRepository(db, image_domain.CollectionImage)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	lesson := &lesson_controller.LessonController{
		LessonUseCase: lesson_usecase.NewLessonUseCase(le, timeout),
		ImageUseCase:  image_usecase.NewImageUseCase(im, timeout),
		AdminUseCase:  admin_usecase.NewAdminUseCase(ad, timeout),
		Database:      env,
	}

	router := group.Group("/lesson")
	router.GET("/fetch", lesson.FetchManyInAdmin)
	router.GET("/fetch/not", lesson.FetchManyNotPaginationInAdmin)
	router.GET("/fetch/_id", lesson.FetchByIdInAdmin)
	router.GET("/fetch/course_id", lesson.FetchByIdCourseInAdmin)
	router.POST("/create", lesson.CreateOneLesson)
	router.POST("/create/0/image", lesson.CreateOneLessonNotImage)
	router.POST("/create/1/image", lesson.CreateOneLessonHaveImage)
	router.POST("/create/file", lesson.CreateLessonWithFile)
	router.PATCH("/update", lesson.UpdateOneLesson)
	router.PATCH("/update/image", lesson.UpdateImageLesson)
	router.DELETE("/delete/:_id", lesson.DeleteOneLesson)
}
