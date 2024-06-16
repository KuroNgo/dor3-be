package vocabulary_route

import (
	vocabulary_controller "clean-architecture/api/controller/vocabulary"
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	exam_domain "clean-architecture/domain/exam"
	exercise_domain "clean-architecture/domain/exercise"
	image_domain "clean-architecture/domain/image"
	lesson_domain "clean-architecture/domain/lesson"
	mark_domain "clean-architecture/domain/mark_vocabulary"
	quiz_domain "clean-architecture/domain/quiz"
	unit_domain "clean-architecture/domain/unit"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
	admin_repository "clean-architecture/repository/admin"
	image_repository "clean-architecture/repository/image"
	unit_repo "clean-architecture/repository/unit"
	vocabulary_repository "clean-architecture/repository/vocabulary"
	admin_usecase "clean-architecture/usecase/admin"
	image_usecase "clean-architecture/usecase/image"
	unit_usecase "clean-architecture/usecase/unit"
	vocabulary_usecase "clean-architecture/usecase/vocabulary"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AdminVocabularyRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	vo := vocabulary_repository.NewVocabularyRepository(db, vocabulary_domain.CollectionVocabulary, mark_domain.CollectionMark, unit_domain.CollectionUnit, lesson_domain.CollectionLesson)
	un := unit_repo.NewUnitRepository(db, unit_domain.CollectionUnit, unit_domain.CollectionUnitProcess, lesson_domain.CollectionLesson,
		vocabulary_domain.CollectionVocabulary, exam_domain.CollectionExam, exercise_domain.CollectionExercise, quiz_domain.CollectionQuiz)
	im := image_repository.NewImageRepository(db, image_domain.CollectionImage)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	vocabulary := &vocabulary_controller.VocabularyController{
		VocabularyUseCase: vocabulary_usecase.NewVocabularyUseCase(vo, timeout),
		UnitUseCase:       unit_usecase.NewUnitUseCase(un, timeout),
		ImageUseCase:      image_usecase.NewImageUseCase(im, timeout),
		AdminUseCase:      admin_usecase.NewAdminUseCase(ad, timeout),
		Database:          env,
	}

	router := group.Group("/vocabulary")
	router.POST("/create", vocabulary.CreateOneVocabulary)
	router.POST("/create/file", vocabulary.CreateVocabularyWithFileInAdmin)
	router.POST("create/audio", vocabulary.GenerateVoice)
	router.PUT("/update/_id", vocabulary.UpdateOneVocabulary)
	router.DELETE("/delete/1/_id", vocabulary.DeleteOneVocabulary)
	router.DELETE("/delete/n/_id", vocabulary.DeleteManyVocabulary)
}
