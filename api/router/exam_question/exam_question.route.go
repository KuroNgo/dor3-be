package exam_question_route

import (
	exam_question_controller "clean-architecture/api/controller/exam_question"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exam_domain "clean-architecture/domain/exam"
	exam_question_domain "clean-architecture/domain/exam_question"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	exam_question_repository "clean-architecture/repository/exam_question"
	user_repository "clean-architecture/repository/user"
	exam_question_usecase "clean-architecture/usecase/exam_question"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExamQuestionRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	quest := exam_question_repository.NewExamQuestionRepository(db, exam_question_domain.CollectionExamQuestion, exam_domain.CollectionExam, vocabulary_domain.CollectionVocabulary)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	question := &exam_question_controller.ExamQuestionsController{
		ExamQuestionUseCase: exam_question_usecase.NewExamQuestionUseCase(quest, timeout),
		UserUseCase:         user_usecase.NewUserUseCase(users, timeout),
		Database:            env,
	}

	router := group.Group("/exam/question")
	router.GET("/fetch", middleware.DeserializeUser(), question.FetchManyExamQuestions)
	router.GET("/fetch/exam_id", middleware.DeserializeUser(), question.FetchManyExamQuestionsByExamID)
}
