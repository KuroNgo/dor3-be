package activity_log_route

import (
	activity_controller "clean-architecture/api/controller/activity"
	"clean-architecture/bootstrap"
	activity_log_domain "clean-architecture/domain/activity_log"
	user_domain "clean-architecture/domain/user"
	"clean-architecture/infrastructor/mongo"
	activity_repository "clean-architecture/repository/activity"
	activity_usecase "clean-architecture/usecase/activity"
	"github.com/gin-gonic/gin"
	"time"
)

func ActivityRoute(env *bootstrap.Database, timeout time.Duration, db mongo.Database) *activity_controller.ActivityController {
	ac := activity_repository.NewActivityRepository(db, activity_log_domain.CollectionActivityLog, user_domain.CollectionUser)
	activity := &activity_controller.ActivityController{
		ActivityUseCase: activity_usecase.NewActivityUseCase(ac, timeout),
		Database:        env,
	}

	return activity
}

func AdminActivityRoute(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ac := activity_repository.NewActivityRepository(db, activity_log_domain.CollectionActivityLog, user_domain.CollectionUser)
	activity := &activity_controller.ActivityController{
		ActivityUseCase: activity_usecase.NewActivityUseCase(ac, timeout),
		Database:        env,
	}

	router := group.Group("/activity")
	router.GET("/fetch", activity.FetchManyActivity)
}
