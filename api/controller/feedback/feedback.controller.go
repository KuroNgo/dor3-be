package feedback_controller

import (
	"clean-architecture/bootstrap"
	admin_domain "clean-architecture/domain/admin"
	feedback_domain "clean-architecture/domain/feedback"
	user_domain "clean-architecture/domain/user"
)

type FeedbackController struct {
	FeedbackUseCase feedback_domain.IFeedbackUseCase
	AdminUseCase    admin_domain.IAdminUseCase
	UserUseCase     user_domain.IUserUseCase
	Database        *bootstrap.Database
}
