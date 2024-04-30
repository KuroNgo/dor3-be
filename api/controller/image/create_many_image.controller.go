package image_controller

import (
	image_domain "clean-architecture/domain/image"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (im *ImageController) CreateManyImageForLesson(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error parsing form",
		})
		return
	}
	files := form.File["files"]

	for _, file := range files {
		if !file_internal.IsImage(file.Filename) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		f, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error opening uploaded file",
			})
			return
		}

		// Tải lên tệp vào Cloudinary
		result, err := cloudinary.UploadImageToCloudinary(f, file.Filename, im.Database.CloudinaryUploadFolderLesson)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error uploading file",
			})
			return
		}

		// Tạo metadata từ thông tin của file và đường dẫn trả về từ Cloudinary
		metadata := &image_domain.Image{
			Id:        primitive.NewObjectID(),
			ImageName: file.Filename,
			ImageUrl:  result.ImageURL,
			Size:      file.Size / 1024,
			Category:  "lesson",
			AssetId:   result.AssetID,
		}

		err = im.ImageUseCase.CreateOne(ctx, metadata)
		if err != nil {
			//ctx.JSON(http.StatusBadRequest, gin.H{
			//	"error": err.Error(),
			//})
			continue
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "Files uploaded successfully",
	})

}

func (im *ImageController) CreateManyImageForExam(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error parsing form",
		})
		return
	}
	files := form.File["files"]

	for _, file := range files {
		if !file_internal.IsImage(file.Filename) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		f, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error opening uploaded file",
			})
			return
		}

		// Tải lên tệp vào Cloudinary
		result, err := cloudinary.UploadImageToCloudinary(f, file.Filename, im.Database.CloudinaryUploadFolderExam)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error uploading file",
			})
			return
		}

		// Tạo metadata từ thông tin của file và đường dẫn trả về từ Cloudinary
		metadata := &image_domain.Image{
			Id:        primitive.NewObjectID(),
			ImageName: file.Filename,
			ImageUrl:  result.ImageURL,
			Size:      file.Size / 1024,
			Category:  "exam",
			AssetId:   result.AssetID,
		}

		err = im.ImageUseCase.CreateOne(ctx, metadata)
		if err != nil {
			//ctx.JSON(http.StatusBadRequest, gin.H{
			//	"error": err.Error(),
			//})
			continue
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "Files uploaded successfully",
	})

}

func (im *ImageController) CreateManyImageForUser(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error parsing form",
		})
		return
	}
	files := form.File["files"]

	for _, file := range files {
		if !file_internal.IsImage(file.Filename) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		f, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error opening uploaded file",
			})
			return
		}

		// Tải lên tệp vào Cloudinary
		result, err := cloudinary.UploadImageToCloudinary(f, file.Filename, im.Database.CloudinaryUploadFolderUser)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error uploading file",
			})
			return
		}

		// Tạo metadata từ thông tin của file và đường dẫn trả về từ Cloudinary
		metadata := &image_domain.Image{
			Id:        primitive.NewObjectID(),
			ImageName: file.Filename,
			ImageUrl:  result.ImageURL,
			Size:      file.Size / 1024,
			Category:  "user",
			AssetId:   result.AssetID,
		}

		err = im.ImageUseCase.CreateOne(ctx, metadata)
		if err != nil {
			//ctx.JSON(http.StatusBadRequest, gin.H{
			//	"error": err.Error(),
			//})
			continue
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "Files uploaded successfully",
	})

}

func (im *ImageController) CreateManyImageForQuiz(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   "Error parsing form",
		})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   "Error parsing form",
		})
		return
	}
	files := form.File["files"]

	for _, file := range files {
		if !file_internal.IsImage(file.Filename) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		f, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error opening uploaded file",
			})
			return
		}

		// Tải lên tệp vào Cloudinary
		result, err := cloudinary.UploadImageToCloudinary(f, file.Filename, im.Database.CloudinaryUploadFolderQuiz)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error uploading file",
			})
			return
		}

		// Tạo metadata từ thông tin của file và đường dẫn trả về từ Cloudinary
		metadata := &image_domain.Image{
			Id:        primitive.NewObjectID(),
			ImageName: file.Filename,
			ImageUrl:  result.ImageURL,
			Size:      file.Size / 1024,
			Category:  "quiz",
			AssetId:   result.AssetID,
		}

		err = im.ImageUseCase.CreateOne(ctx, metadata)
		if err != nil {
			//ctx.JSON(http.StatusBadRequest, gin.H{
			//	"error": err.Error(),
			//})
			continue
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "Files uploaded successfully",
	})

}

func (im *ImageController) CreateManyImageForStatic(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   "Error parsing form",
		})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error parsing form",
		})
		return
	}
	files := form.File["files"]

	for _, file := range files {
		if !file_internal.IsImage(file.Filename) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		f, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error opening uploaded file",
			})
			return
		}

		// Tải lên tệp vào Cloudinary
		result, err := cloudinary.UploadImageToCloudinary(f, file.Filename, im.Database.CloudinaryUploadFolderStatic)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error uploading file",
			})
			return
		}

		// Tạo metadata từ thông tin của file và đường dẫn trả về từ Cloudinary
		metadata := &image_domain.Image{
			Id:        primitive.NewObjectID(),
			ImageName: file.Filename,
			ImageUrl:  result.ImageURL,
			Size:      file.Size / 1024,
			Category:  "static",
			AssetId:   result.AssetID,
		}

		err = im.ImageUseCase.CreateOne(ctx, metadata)
		if err != nil {
			//ctx.JSON(http.StatusBadRequest, gin.H{
			//	"error": err.Error(),
			//})
			continue
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "Files uploaded successfully",
	})

}
