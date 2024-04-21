package image_route

import (
	image_controller "clean-architecture/api/controller/image"
	"clean-architecture/bootstrap"
	image_domain "clean-architecture/domain/image"
	image_repository "clean-architecture/repository/image"
	image_usecase "clean-architecture/usecase/image"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ImageRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	im := image_repository.NewImageRepository(db, image_domain.CollectionImage)
	image := &image_controller.ImageController{
		ImageUseCase: image_usecase.NewImageUseCase(im, timeout),
		Database:     env,
	}

	router := group.Group("/image")
	router.GET("/fetch-by-name", image.FetchImageByName)
	router.GET("/fetch-many", image.FetchImage)
}
