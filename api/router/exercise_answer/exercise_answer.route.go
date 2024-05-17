package exercise_answer_route

import (
	exercise_answer_controller "clean-architecture/api/controller/exercise_answer"
	"clean-architecture/api/middleware"
	"clean-architecture/bootstrap"
	exercise_answer_domain "clean-architecture/domain/exercise_answer"
	exercise_questions_domain "clean-architecture/domain/exercise_questions"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	exercise_answer_repository "clean-architecture/repository/exercise_answer"
	user_repository "clean-architecture/repository/user"
	exercise_answer_usecase "clean-architecture/usecase/exercise_answer"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ExerciseRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ans := exercise_answer_repository.NewExerciseAnswerRepository(db, exercise_questions_domain.CollectionExerciseQuestion, exercise_answer_domain.CollectionExerciseAnswers)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	answer := &exercise_answer_controller.ExerciseAnswerController{
		ExerciseAnswerUseCase: exercise_answer_usecase.NewExerciseAnswerUseCase(ans, timeout),
		UserUseCase:           usecase.NewUserUseCase(ur, timeout),
		Database:              env,
	}

	router := group.Group("/exercise/answer")
	router.Use(middleware.DeserializeUser())
	router.GET("/fetch", answer.FetchManyAnswerByUserIDAndQuestionID)
	router.POST("/create", answer.CreateOneExerciseAnswer)
	router.DELETE("/delete", answer.DeleteOneAnswer)
	router.GET("/all/delete", answer.DeleteAllAnswerInExerciseID)
}
