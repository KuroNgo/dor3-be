package exercise_question_route

import (
	exercise_quesiton_controller "clean-architecture/api/controller/exercise_quesiton"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exercise_domain "clean-architecture/domain/exercise"
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	exercise_question_repository "clean-architecture/repository/exercise_question"
	user_repository "clean-architecture/repository/user"
	exercise_question_usecase "clean-architecture/usecase/exercise_question"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExerciseQuestionRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	quest := exercise_question_repository.NewExerciseQuestionRepository(db, exercise_questions_domain.CollectionExerciseQuestion, exercise_domain.CollectionExercise)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	question := &exercise_quesiton_controller.ExerciseQuestionsController{
		ExerciseQuestionUseCase: exercise_question_usecase.NewExerciseQuestionUseCase(quest, timeout),
		UserUseCase:             user_usecase.NewUserUseCase(users, timeout),
		Database:                env,
	}

	router := group.Group("/exercise/question")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch/_id", question.FetchOneExerciseQuestionByID)
	router.GET("/fetch", question.FetchManyExerciseOptions)
}
