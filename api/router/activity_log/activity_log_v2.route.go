package activity_log_route

import (
	activity_controller "clean-architecture/api/controller/activity"
	"clean-architecture/bootstrap"
	activity_log_domain "clean-architecture/domain/activity_log"
	activity_repository "clean-architecture/repository/activity"
	activity_usecase "clean-architecture/usecase/activity"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ActivityRouteV2(env *bootstrap.Database, timeout time.Duration, db mongo.Database) *activity_controller.ActivityControllerV2 {
	ac := activity_repository.NewActivityRepositoryV2(db, activity_log_domain.CollectionActivityLog)
	activity := &activity_controller.ActivityControllerV2{
		ActivityUseCase: activity_usecase.NewActivityUseCaseV2(ac, timeout),
		Database:        env,
	}

	return activity
}
