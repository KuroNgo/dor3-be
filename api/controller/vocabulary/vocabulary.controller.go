package vocabulary_controller

import (
	"clean-architecture/bootstrap"
	vocabulary_domain "clean-architecture/domain/vocabulary"
)

type VocabularyController struct {
	VocabularyUseCase vocabulary_domain.IVocabularyUseCase
	Database          *bootstrap.Database
}
