package vocabulary_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
)

type VocabularyController struct {
	VocabularyUseCase vocabulary_domain.IVocabularyUseCase
	UserUseCase       user_domain.IUserUseCase
	AdminUseCase      admin_domain.IAdminUseCase
	Database          *bootstrap.Database
}
