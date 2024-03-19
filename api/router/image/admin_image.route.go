package image_route

import (
	image_controller "clean-architecture/api/controller/image"
	"clean-architecture/bootstrap"
	image_domain "clean-architecture/domain/image"
	"clean-architecture/infrastructor/mongo"
	"clean-architecture/internal/cloud/cloudinary"
	image_repository "clean-architecture/repository/image"
	image_usecase "clean-architecture/usecase/image"
	"github.com/gin-gonic/gin"
	"time"
)

func AdminImageRoute(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	im := image_repository.NewImageRepository(db, image_domain.CollectionImage)
	image := &image_controller.ImageController{
		ImageUseCase: image_usecase.NewImageUseCase(im, timeout),
		Database:     env,
	}

	router := group.Group("/image")
	// riêng audio thì giành cho cả 2 phiá
	router.POST("/file", cloudinary.FileUploadMiddleware(), image.CreateImageInCloudinaryAndSaveMetadataInDatabase)
}
