package user_controller

import (
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (u *UserController) UpdateUser(ctx *gin.Context) {
	cookie, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	sub, err := internal.ValidateToken(cookie, u.Database.AccessTokenPublicKey)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	var userInput user_domain.Input
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	result, err := u.UserUseCase.GetByID(ctx, fmt.Sprint(sub))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		userResponse := user_domain.User{
			FullName:   result.FullName,
			Email:      userInput.Email,
			Phone:      userInput.Phone,
			Role:       result.Role,
			UpdatedAt:  time.Now(),
			Specialize: userInput.Specialize,
		}

		_, err = u.UserUseCase.Update(ctx, &userResponse)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
		return
	}

	// Nếu có file được tải lên, tiến hành xử lý file
	file, err = ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	// Kiểm tra xem file có phải là hình ảnh không
	if !file_internal.IsImage(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file format. Only images are allowed.",
		})
		return
	}

	// Mở file để đọc dữ liệu
	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error opening uploaded file",
		})
		return
	}
	defer f.Close()

	// Tải file lên Cloudinary
	filename, _ := ctx.Get("filePath")
	imageURL, err := cloudinary.UploadToCloudinary(f, filename.(string), u.Database.CloudinaryUploadFolderUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resultString, err := json.Marshal(result)
	userResponse := user_domain.User{
		FullName:   result.FullName,
		Email:      userInput.Email,
		Phone:      userInput.Phone,
		Role:       result.Role,
		AvatarURL:  imageURL.ImageURL,
		AssetID:    imageURL.AssetID,
		UpdatedAt:  time.Now(),
		Specialize: userInput.Specialize,
	}

	_, err = u.UserUseCase.Update(ctx, &userResponse)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": string(resultString) + "the user belonging to this token no logger exists",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Updated user",
	})

}
