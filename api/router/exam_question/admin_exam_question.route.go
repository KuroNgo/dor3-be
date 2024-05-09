package exam_question_route

import (
	exam_question_controller "clean-architecture/api/controller/exam_question"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exam_domain "clean-architecture/domain/exam"
	exam_question_domain "clean-architecture/domain/exam_question"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	admin_repository "clean-architecture/repository/admin"
	"clean-architecture/repository/exam_question"
	admin_usecase "clean-architecture/usecase/admin"
	exam_question_usecase "clean-architecture/usecase/exam_question"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminExamQuestionRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	quest := exam_question_repository.NewExamQuestionRepository(db, exam_question_domain.CollectionExamQuestion, exam_domain.CollectionExam, vocabulary_domain.CollectionVocabulary)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	question := &exam_question_controller.ExamQuestionsController{
		ExamQuestionUseCase: exam_question_usecase.NewExamQuestionUseCase(quest, timeout),
		AdminUseCase:        admin_usecase.NewAdminUseCase(ad, timeout),
		Database:            env,
	}

	router := group.Group("/exam/question")
	router.POST("/create", question.CreateOneExamQuestions)
	router.PATCH("/update", question.UpdateOneExamQuestion)
	router.DELETE("/delete/_id", question.DeleteOneExamQuestions)
}
