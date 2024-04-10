package mark_vocabulary_controller

import (
	"clean-architecture/bootstrap"
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	user_domain "clean-architecture/domain/user"
)

type MarkVocabularyController struct {
	MarkVocabularyUseCase mark_vocabulary_domain.IMarkToFavouriteUseCase
	UserUseCase           user_domain.IUserUseCase
	Database              *bootstrap.Database
}
