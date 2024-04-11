package mark_vocabulary_controller

import (
	"clean-architecture/bootstrap"
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	user_domain "clean-architecture/domain/user"
	vocabulary_domain "clean-architecture/domain/vocabulary"
)

type MarkVocabularyController struct {
	MarkVocabularyUseCase mark_vocabulary_domain.IMarkToFavouriteUseCase
	VocabularyUseCase     vocabulary_domain.IVocabularyUseCase
	UserUseCase           user_domain.IUserUseCase
	Database              *bootstrap.Database
}
