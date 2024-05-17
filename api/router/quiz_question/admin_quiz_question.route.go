package quiz_question_route

import (
	quiz_question_controller "clean-architecture/api/controller/quiz_question"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	quiz_domain "clean-architecture/domain/quiz"
	quiz_question_domain "clean-architecture/domain/quiz_question"
	user_domain "clean-architecture/domain/user"
	admin_repository "clean-architecture/repository/admin"
	quiz_question_repository "clean-architecture/repository/quiz_question"
	admin_usecase "clean-architecture/usecase/admin"
	quiz_question_usecase "clean-architecture/usecase/quiz_question"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminQuizQuestionRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	quest := quiz_question_repository.NewQuizQuestionRepository(db, quiz_question_domain.CollectionQuizQuestion, quiz_domain.CollectionQuiz)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	question := &quiz_question_controller.QuizQuestionsController{
		QuizQuestionUseCase: quiz_question_usecase.NewQuizQuestionUseCase(quest, timeout),
		AdminUseCase:        admin_usecase.NewAdminUseCase(ad, timeout),
		Database:            env,
	}

	router := group.Group("/quiz/question")
	router.GET("/fetch/_id", question.FetchOneQuizQuestionByID)
	router.POST("/create", question.CreateOneQuizQuestions)
	router.PATCH("/update", question.UpdateOneQuizOptions)
	router.DELETE("/delete/_id", question.DeleteOneQuizQuestions)
}
