package user_controller

import (
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (u *UserController) UpdateUser(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	user, err := u.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	fullName := ctx.Request.FormValue("full_name")
	phone := ctx.Request.FormValue("phone")

	file, err := ctx.FormFile("file")
	if err != nil {
		userResponse := user_domain.User{
			ID:        user.ID,
			FullName:  fullName,
			Phone:     phone,
			Role:      user.Role,
			UpdatedAt: time.Now(),
		}

		err = u.UserUseCase.Update(ctx, &userResponse)
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
	imageURL, err := cloudinary.UploadImageToCloudinary(f, filename.(string), u.Database.CloudinaryUploadFolderUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resultString, err := json.Marshal(user)
	userResponse := user_domain.User{
		FullName:  fullName,
		Phone:     phone,
		Role:      user.Role,
		AvatarURL: imageURL.ImageURL,
		AssetID:   imageURL.AssetID,
		UpdatedAt: time.Now(),
	}

	err = u.UserUseCase.Update(ctx, &userResponse)
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
