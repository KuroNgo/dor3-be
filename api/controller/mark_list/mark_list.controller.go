package marklist_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	mark_list_domain "clean-architecture/domain/mark_list"
	mark_vocabulary_domain "clean-architecture/domain/mark_vocabulary"
	user_domain "clean-architecture/domain/user"
)

type MarkListController struct {
	MarkListUseCase       mark_list_domain.IMarkListUseCase
	MarkVocabularyUseCase mark_vocabulary_domain.IMarkToFavouriteUseCase
	AdminUseCase          admin_domain.IAdminUseCase
	UserUseCase           user_domain.IUserUseCase
	Database              *bootstrap.Database
}
