package vocabulary_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	image_domain "clean-architecture/domain/image"
	vocabulary_domain "clean-architecture/domain/vocabulary"
)

type VocabularyController struct {
	VocabularyUseCase vocabulary_domain.IVocabularyUseCase
	ImageUseCase      image_domain.IImageUseCase
	AdminUseCase      admin_domain.IAdminUseCase
	Database          *bootstrap.Database
}
