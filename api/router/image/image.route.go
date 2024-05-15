package image_route

import (
	image_controller "clean-architecture/api/controller/image"
	"clean-architecture/bootstrap"
	image_domain "clean-architecture/domain/image"
	user_domain "clean-architecture/domain/user"
	user_detail_domain "clean-architecture/domain/user_detail"
	image_repository "clean-architecture/repository/image"
	user_repository "clean-architecture/repository/user"
	image_usecase "clean-architecture/usecase/image"
	usecase "clean-architecture/usecase/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ImageRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	im := image_repository.NewImageRepository(db, image_domain.CollectionImage)
	ur := user_repository.NewUserRepository(db, user_domain.CollectionUser, user_detail_domain.CollectionUserDetail)

	image := &image_controller.ImageController{
		ImageUseCase: image_usecase.NewImageUseCase(im, timeout),
		UserUseCase:  usecase.NewUserUseCase(ur, timeout),
		Database:     env,
	}

	router := group.Group("/image")
	router.GET("/fetch/name", image.FetchImageByName)
	router.GET("/fetch/category", image.FetchImageByCategory)
	router.GET("/fetch", image.FetchImage)
}
