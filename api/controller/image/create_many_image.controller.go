package image_controller

import (
	image_domain "clean-architecture/domain/image"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"net/http"
)

func (im *ImageController) CreateManyImageForLesson(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error parsing form",
		})
		return
	}

	files := form.File["files"]
	errChan := make(chan error)

	for _, file := range files {
		go func(file *multipart.FileHeader) {
			if !file_internal.IsImage(file.Filename) {
				errChan <- errors.New("invalid file format")
				return
			}

			f, err := file.Open()
			if err != nil {
				errChan <- err
				return
			}

			// Tải lên tệp vào Cloudinary
			result, err := cloudinary.UploadToCloudinary(f, file.Filename, im.Database.CloudinaryUploadFolderLesson)
			if err != nil {
				errChan <- err
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

			// Gọi UseCase để lưu hình ảnh vào cơ sở dữ liệu
			err = im.ImageUseCase.CreateOne(ctx, metadata)
			if err != nil {
				errChan <- err
				return
			}

			// Gửi một thông báo thành công đến kênh
			errChan <- nil
		}(file)
	}

	// Xử lý các lỗi trả về từ các Goroutine
	for range files {
		if err := <-errChan; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error processing file",
			})
			return
		}
	}

	// Trả về thông báo thành công nếu không có lỗi
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"success": "Files uploaded successfully",
	})
}

func (im *ImageController) CreateManyImageForExam(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error parsing form",
		})
		return
	}
	files := form.File["files"]
	errChan := make(chan error)

	for _, file := range files {
		go func(file *multipart.FileHeader) {
			if !file_internal.IsImage(file.Filename) {
				errChan <- errors.New("invalid file format")
				return
			}

			f, err := file.Open()
			if err != nil {
				errChan <- err
				return
			}

			// Tải lên tệp vào Cloudinary
			result, err := cloudinary.UploadToCloudinary(f, file.Filename, im.Database.CloudinaryUploadFolderLesson)
			if err != nil {
				errChan <- err
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
				errChan <- err
				return
			}
			// Gửi một thông báo thành công đến kênh
			errChan <- nil
		}(file)
	}

	// Xử lý các lỗi trả về từ các Goroutine
	for range files {
		if err := <-errChan; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error processing file",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "Files uploaded successfully",
	})

}

func (im *ImageController) CreateManyImageForUser(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error parsing form",
		})
		return
	}
	files := form.File["files"]
	errChan := make(chan error)

	for _, file := range files {
		go func(file *multipart.FileHeader) {
			if !file_internal.IsImage(file.Filename) {
				errChan <- err
				return
			}

			f, err := file.Open()
			if err != nil {
				errChan <- err
				return
			}

			// Tải lên tệp vào Cloudinary
			result, err := cloudinary.UploadToCloudinary(f, file.Filename, im.Database.CloudinaryUploadFolderUser)
			if err != nil {
				errChan <- err
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
				errChan <- err
			}
		}(file)
	}

	// Xử lý các lỗi trả về từ các Goroutine
	for range files {
		if err := <-errChan; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error processing file",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "Files uploaded successfully",
	})

}

func (im *ImageController) CreateManyImageForQuiz(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10MB max size
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   "Error parsing form",
		})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   "Error parsing form",
		})
		return
	}

	files := form.File["files"]
	errChan := make(chan error)

	for _, file := range files {
		go func(file *multipart.FileHeader) {
			if !file_internal.IsImage(file.Filename) {
				errChan <- err
				return
			}

			f, err := file.Open()
			if err != nil {
				errChan <- err
				return
			}

			// Tải lên tệp vào Cloudinary
			result, err := cloudinary.UploadToCloudinary(f, file.Filename, im.Database.CloudinaryUploadFolderQuiz)
			if err != nil {
				errChan <- err
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
				errChan <- err
			}
		}(file)
	}

	// Xử lý các lỗi trả về từ các Goroutine
	for range files {
		if err := <-errChan; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error processing file",
			})
			return
		}
	}
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"success": "Files uploaded successfully",
	})

}

func (im *ImageController) CreateManyImageForStatic(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   "Error parsing form",
		})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error parsing form",
		})
		return
	}
	files := form.File["files"]
	errChan := make(chan error)
	for _, file := range files {
		go func(file *multipart.FileHeader) {
			if !file_internal.IsImage(file.Filename) {
				errChan <- err
				return
			}

			f, err := file.Open()
			if err != nil {
				errChan <- err
				return
			}

			// Tải lên tệp vào Cloudinary
			result, err := cloudinary.UploadToCloudinary(f, file.Filename, im.Database.CloudinaryUploadFolderStatic)
			if err != nil {
				errChan <- err
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
				errChan <- err
			}
		}(file)
	}

	// Xử lý các lỗi trả về từ các Goroutine
	for range files {
		if err := <-errChan; err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error processing file",
			})
			return
		}
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"success": "Files uploaded successfully",
	})

}
