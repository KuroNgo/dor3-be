package exercise_answer_route

import (
	exercise_answer_controller "clean-architecture/api/controller/exercise_answer"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	quiz_answer_domain "clean-architecture/domain/quiz_answer"
	quiz_question_domain "clean-architecture/domain/quiz_question"
	user_domain "clean-architecture/domain/user"
	exercise_answer_repository "clean-architecture/repository/exercise_answer"
	user_repository "clean-architecture/repository/user"
	exercise_answer_usecase "clean-architecture/usecase/exercise_answer"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExerciseRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ans := exercise_answer_repository.NewExerciseAnswerRepository(db, quiz_question_domain.CollectionQuizQuestion, quiz_answer_domain.CollectionQuizAnswers)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser)

	answer := &exercise_answer_controller.ExerciseAnswerController{
		ExerciseAnswerUseCase: exercise_answer_usecase.NewExerciseAnswerUseCase(ans, timeout),
		UserUseCase:           usecase.NewUserUseCase(ur, timeout),
		Database:              env,
	}

	router := group.Group("/exercise/answer")
	router.GET("/fetch", middleware.DeserializeUser(), answer.FetchManyAnswerByUserIDAndQuestionID)
	router.POST("/create", middleware.DeserializeUser(), answer.CreateOneExerciseAnswer)
	router.DELETE("/delete", middleware.DeserializeUser(), answer.DeleteOneAnswer)
	router.GET("/delete/all", middleware.DeserializeUser(), answer.DeleteAllAnswerInExerciseID)
}
