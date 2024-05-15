package activity_log_route

import (
	activity_controller "clean-architecture/api/controller/activity"
	"clean-architecture/bootstrap"
	activity_log_domain "clean-architecture/domain/activity_log"
	admin_domain "clean-architecture/domain/admin"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	activity_repository "clean-architecture/repository/activity"
	admin_repository "clean-architecture/repository/admin"
	user_repository "clean-architecture/repository/user"
	activity_usecase "clean-architecture/usecase/activity"
	admin_usecase "clean-architecture/usecase/admin"
	user_usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ActivityRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database) *activity_controller.ActivityController {
	ac := activity_repository.NewActivityRepository(db, activity_log_domain.CollectionActivityLog)
	users := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	activity := &activity_controller.ActivityController{
		ActivityUseCase: activity_usecase.NewActivityUseCase(ac, timeout),
		UserUseCase:     user_usecase.NewUserUseCase(users, timeout),
		AdminUseCase:    admin_usecase.NewAdminUseCase(ad, timeout),
		Database:        env,
	}

	return activity
}

func AdminActivityRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	ac := activity_repository.NewActivityRepository(db, activity_log_domain.CollectionActivityLog)
	ad := admin_repository.NewAdminRepository(db, admin_domain.CollectionAdmin, user_domain.CollectionUser)

	activity := &activity_controller.ActivityController{
		ActivityUseCase: activity_usecase.NewActivityUseCase(ac, timeout),
		AdminUseCase:    admin_usecase.NewAdminUseCase(ad, timeout),
		Database:        env,
	}

	router := group.Group("/activity")
	router.GET("/fetch", activity.FetchManyActivity)
}
