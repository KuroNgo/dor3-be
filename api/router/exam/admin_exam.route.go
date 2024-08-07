package exam_route

import (
	exam_controller "clean-architecture/api/controller/exam"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exam_domain "clean-architecture/domain/exam"
	exam_question_domain "clean-architecture/domain/exam_question"
	lesson_domain "clean-architecture/domain/lesson"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	admin_repository "clean-architecture/repository/admin"
	exam_repository "clean-architecture/repository/exam"
	admin_usecase "clean-architecture/usecase/admin"
	exam_usecase "clean-architecture/usecase/exam"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminExamRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ex := exam_repository.NewExamRepository(db, exam_domain.CollectionExam, exam_domain.CollectionExamProcess, lesson_domain.CollectionLesson, unit_domain.CollectionUnit, exam_question_domain.CollectionExamQuestion, vocabulary_domain.CollectionVocabulary)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	exam := &exam_controller.ExamsController{
		ExamUseCase:  exam_usecase.NewExamUseCase(ex, timeout),
		AdminUseCase: admin_usecase.NewAdminUseCase(ad, timeout),
		Database:     env,
	}

	router := group.Group("/exam")
	router.GET("/fetch/_id", exam.FetchOneExamByIDInAdmin)
	router.GET("/all/fetch", exam.FetchManyExamInAdmin)
	router.GET("/fetch/ns/unit_id", exam.FetchManyExamByUnitIDInAdmin)
	router.GET("/fetch/1s/unit_id", exam.FetchOneExamByUnitIDInAdmin)
	router.POST("/create", exam.CreateOneExamInAdmin)
	router.PATCH("/update", exam.UpdateOneExamInAdmin)
	router.DELETE("/delete/_id", exam.DeleteOneExamInAdmin)
}
