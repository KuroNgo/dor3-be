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

	router.POST("/file/upload/1/static", cloudinary.FileUploadMiddleware(), image.CreateOneImageStatic)
	router.POST("/file/upload/1/lesson", cloudinary.FileUploadMiddleware(), image.CreateOneImageLesson)
	router.POST("/file/upload/1/exam", cloudinary.FileUploadMiddleware(), image.CreateOneImageExam)
	router.POST("/file/upload/1/user", cloudinary.FileUploadMiddleware(), image.CreateOneImageUser)
	router.POST("/file/upload/1/quiz", cloudinary.FileUploadMiddleware(), image.CreateOneImageQuiz)

	router.POST("/files/upload/many/user", cloudinary.FileUploadMiddleware(), image.CreateManyImageForUser)
	router.POST("/files/upload/many/quiz", cloudinary.FileUploadMiddleware(), image.CreateManyImageForQuiz)
	router.POST("/files/upload/many/static", cloudinary.FileUploadMiddleware(), image.CreateManyImageForStatic)
	router.POST("/files/upload/many/lesson", cloudinary.FileUploadMiddleware(), image.CreateManyImageForLesson)
	router.POST("/files/upload/many/exam", cloudinary.FileUploadMiddleware(), image.CreateManyImageForExam)
}
