package course_route

import (
	course_controller "clean-architecture/api/controller/course"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	course_domain "clean-architecture/domain/course"
	exam_domain "clean-architecture/domain/exam"
	exercise_domain "clean-architecture/domain/exercise"
	lesson_domain "clean-architecture/domain/lesson"
	mark_domain "clean-architecture/domain/mark_vocabulary"
	quiz_domain "clean-architecture/domain/quiz"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	admin_repository "clean-architecture/repository/admin"
	course_repository "clean-architecture/repository/course"
	lesson_repository "clean-architecture/repository/lesson"
	unit_repo "clean-architecture/repository/unit"
	vocabulary_repository "clean-architecture/repository/vocabulary"
	admin_usecase "clean-architecture/usecase/admin"
	course_usecase "clean-architecture/usecase/course"
	lesson_usecase "clean-architecture/usecase/lesson"
	unit_usecase "clean-architecture/usecase/unit"
	vocabulary_usecase "clean-architecture/usecase/vocabulary"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminCourseRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	co := course_repository.NewCourseRepository(db, course_domain.CollectionCourse, course_domain.CollectionCourseProcess, lesson_domain.CollectionLesson, unit_domain.CollectionUnit, vocabulary_domain.CollectionVocabulary)
	le := lesson_repository.NewLessonRepository(db, lesson_domain.CollectionLesson, lesson_domain.CollectionLessonProcess, course_domain.CollectionCourse, unit_domain.CollectionUnit, vocabulary_domain.CollectionVocabulary)
	un := unit_repo.NewUnitRepository(db, unit_domain.CollectionUnit, lesson_domain.CollectionLesson, vocabulary_domain.CollectionVocabulary, exam_domain.CollectionExam, exercise_domain.CollectionExercise, quiz_domain.CollectionQuiz)
	vo := vocabulary_repository.NewVocabularyRepository(db, vocabulary_domain.CollectionVocabulary, mark_domain.CollectionMark, unit_domain.CollectionUnit, lesson_domain.CollectionLesson)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	course := &course_controller.CourseController{
		CourseUseCase:     course_usecase.NewCourseUseCase(co, timeout),
		LessonUseCase:     lesson_usecase.NewLessonUseCase(le, timeout),
		UnitUseCase:       unit_usecase.NewUnitUseCase(un, timeout),
		VocabularyUseCase: vocabulary_usecase.NewVocabularyUseCase(vo, timeout),
		AdminUseCase:      admin_usecase.NewAdminUseCase(ad, timeout),
		Database:          env,
	}

	router := group.Group("/course")
	router.GET("/fetch", course.FetchCourseInAdmin)
	router.GET("/fetch/_id", course.FetchCourseByIDInAdmin)
	router.POST("/create", course.CreateOneCourseInAdmin)
	router.POST("/create/file", course.CreateCourseWithFileInAdmin)
	router.POST("/create/file/final", course.CreateLessonManagementWithFileInAdmin)
	router.PATCH("/update", course.UpdateCourseInAdmin)
	router.DELETE("/delete/:_id", course.DeleteCourseInAdmin)
}
