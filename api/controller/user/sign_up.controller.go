package user_controller

import (
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal"
	"clean-architecture/internal/cloud/cloudinary"
	"clean-architecture/internal/cloud/google"
	subject_const "clean-architecture/internal/cloud/google/const"
	file_internal "clean-architecture/internal/file"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (u *UserController) SignUp(ctx *gin.Context) {
	var user user_domain.SignUp
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	if !internal.EmailValid(user.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email Invalid !",
		})
		return
	}

	// Bên phía client sẽ phải so sánh password thêm một lần nữa đã đúng chưa
	if !internal.PasswordStrong(user.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": "Password must have at least 8 characters, " +
				"including uppercase letters, lowercase letters and numbers!",
		})
		return
	}

	// Băm mật khẩu
	hashedPassword, err := internal.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	user.Password = hashedPassword
	user.Password = internal.Santize(user.Password)
	user.Email = internal.Santize(user.Email)
	file, err := ctx.FormFile("file")

	if err != nil {
		newUser := user_domain.User{
			ID:       primitive.NewObjectID(),
			FullName: user.FullName,
			//AvatarURL:  user.AvatarURL,
			Specialize: user.Specialize,
			Email:      user.Email,
			Password:   hashedPassword,
			Verified:   false,
			Provider:   "fe-it",
			Role:       "user",
			Phone:      user.Phone,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		// thực hiện đăng ký người dùng
		err = u.UserUseCase.Create(ctx, newUser)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error()},
			)
			return
		}

		// Thêm công việc cron để gửi email nhắc nhở
		err = google.SendEmail(user.Email, subject_const.SignInTheFirstTime, subject_const.ContentTitle2)
		if err != nil {
			return
		}

		// Trả về phản hồi thành công
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
	newUser := user_domain.User{
		ID:         primitive.NewObjectID(),
		FullName:   user.FullName,
		AvatarURL:  imageURL.ImageURL,
		AssetID:    imageURL.AssetID,
		Specialize: user.Specialize,
		Email:      user.Email,
		Password:   hashedPassword,
		Verified:   false,
		Provider:   "fe-it",
		Role:       "user",
		Phone:      user.Phone,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// thực hiện đăng ký người dùng
	err = u.UserUseCase.Create(ctx, newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	// add cron job
	err = google.Cron.AddFunc("0h0m1s", func() {
		err = google.SendEmail(user.Email, subject_const.SignInTheFirstTime, subject_const.ContentTitle2)
		if err != nil {
			return
		}
	})
	google.Cron.Start()

	if err != nil {
		return
	}

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
	return
}
