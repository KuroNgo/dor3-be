package course_route

import (
	course_controller "clean-architecture/api/controller/course"
	"clean-architecture/bootstrap"
	course_domain "clean-architecture/domain/course"
	lesson_domain "clean-architecture/domain/lesson"
	mark_domain "clean-architecture/domain/mark_vocabulary"
	mean_domain "clean-architecture/domain/mean"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	course_repository "clean-architecture/repository/course"
	lesson_repository "clean-architecture/repository/lesson"
	unit_repo "clean-architecture/repository/unit"
	user_repository "clean-architecture/repository/user"
	vocabulary_repository "clean-architecture/repository/vocabulary"
	course_usecase "clean-architecture/usecase/course"
	lesson_usecase "clean-architecture/usecase/lesson"
	unit_usecase "clean-architecture/usecase/unit"
	usecase "clean-architecture/usecase/user"
	vocabulary_usecase "clean-architecture/usecase/vocabulary"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminCourseRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	co := course_repository.NewCourseRepository(db, course_domain.CollectionCourse, lesson_domain.CollectionLesson)
	le := lesson_repository.NewLessonRepository(db, lesson_domain.CollectionLesson, course_domain.CollectionCourse, unit_domain.CollectionUnit)
	un := unit_repo.NewUnitRepository(db, unit_domain.CollectionUnit, lesson_domain.CollectionLesson, vocabulary_domain.CollectionVocabulary)
	vo := vocabulary_repository.NewVocabularyRepository(db, vocabulary_domain.CollectionVocabulary, mean_domain.CollectionMean, mark_domain.CollectionMark, unit_domain.CollectionUnit)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	course := &course_controller.CourseController{
		CourseUseCase:     course_usecase.NewCourseUseCase(co, timeout),
		LessonUseCase:     lesson_usecase.NewLessonUseCase(le, timeout),
		UnitUseCase:       unit_usecase.NewUnitUseCase(un, timeout),
		UserUseCase:       usecase.NewUserUseCase(ur, timeout),
		VocabularyUseCase: vocabulary_usecase.NewVocabularyUseCase(vo, timeout),
		Database:          env,
	}

	router := group.Group("/course")
	router.POST("/create", course.CreateOneCourse)
	router.POST("/create/file", course.CreateCourseWithFile)
	router.POST("/create/file/final", course.CreateLessonManagementWithFile)
	router.PUT("/update/:_id", course.UpdateCourse)
	router.POST("/upsert/:_id", course.UpsertOneQuiz)
	router.DELETE("/delete/:_id", course.DeleteCourse)
}
