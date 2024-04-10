package marklist_controller

import (
	"clean-architecture/bootstrap"
	mark_list_domain "clean-architecture/domain/mark_list"
	user_domain "clean-architecture/domain/user"
)

type MarkListController struct {
	MarkListUseCase mark_list_domain.IMarkListUseCase
	UserUseCase     user_domain.IUserUseCase
	Database        *bootstrap.Database
}
