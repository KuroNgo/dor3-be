package exercise_question_route

import (
	exercise_quesiton_controller "clean-architecture/api/controller/exercise_quesiton"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exercise_domain "clean-architecture/domain/exercise"
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
	user_domain "clean-architecture/domain/user"
	admin_repository "clean-architecture/repository/admin"
	"clean-architecture/repository/exercise_question"
	admin_usecase "clean-architecture/usecase/admin"
	exercise_question_usecase "clean-architecture/usecase/exercise_question"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminExerciseQuestionRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	quest := exercise_question_repository.NewExerciseQuestionRepository(db, exercise_questions_domain.CollectionExerciseQuestion, exercise_domain.CollectionExercise)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	question := &exercise_quesiton_controller.ExerciseQuestionsController{
		ExerciseQuestionUseCase: exercise_question_usecase.NewExerciseQuestionUseCase(quest, timeout),
		AdminUseCase:            admin_usecase.NewAdminUseCase(ad, timeout),
		Database:                env,
	}

	router := group.Group("/exercise/question")
	router.GET("/fetch/_id", question.FetchOneExerciseQuestionByID)
	router.POST("/create", question.CreateOneExerciseQuestions)
	router.PATCH("/update", question.UpdateOneExerciseOptions)
	router.DELETE("/delete/:_id", question.DeleteOneExerciseQuestions)
}
