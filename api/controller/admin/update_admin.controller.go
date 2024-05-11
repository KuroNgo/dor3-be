package admin_controller

import (
	admin_domain "clean-architecture/domain/admin"
	"clean-architecture/internal"
	"clean-architecture/internal/cloud/cloudinary"
	file_internal "clean-architecture/internal/file"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (a *AdminController) UpdateAdmin(ctx *gin.Context) {
	cookie, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	sub, err := internal.ValidateToken(cookie, a.Database.AccessTokenPublicKey)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	phone := ctx.Request.FormValue("phone")

	result, err := a.AdminUseCase.GetByID(ctx, fmt.Sprint(sub))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		adminResponse := admin_domain.Admin{
			FullName:  result.FullName,
			Email:     result.Email,
			Phone:     phone,
			Role:      result.Role,
			UpdatedAt: time.Now(),
		}

		_, err = a.AdminUseCase.UpdateOne(ctx, &adminResponse)
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
	imageURL, err := cloudinary.UploadImageToCloudinary(f, filename.(string), a.Database.CloudinaryUploadFolderUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resultString, err := json.Marshal(result)
	adminResponse := admin_domain.Admin{
		FullName:  result.FullName,
		Email:     result.Email,
		Phone:     result.Phone,
		Role:      result.Role,
		AvatarURL: imageURL.ImageURL,
		AssetURL:  imageURL.AssetID,
		UpdatedAt: time.Now(),
	}

	_, err = a.AdminUseCase.UpdateOne(ctx, &adminResponse)
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
