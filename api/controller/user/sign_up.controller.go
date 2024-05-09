package user_controller

import (
	user_domain "clean-architecture/domain/user"
	"clean-architecture/internal"
	"clean-architecture/internal/cloud/cloudinary"
	"clean-architecture/internal/cloud/google"
	file_internal "clean-architecture/internal/file"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (u *UserController) SignUp(ctx *gin.Context) {
	email := ctx.Request.FormValue("email")
	fullName := ctx.Request.FormValue("full_name")
	password := ctx.Request.FormValue("password")
	avatarUrl := ctx.Request.FormValue("avatar_url")
	phone := ctx.Request.FormValue("phone")

	if !internal.EmailValid(email) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email Invalid !",
		})
		return
	}

	// Bên phía client sẽ phải so sánh password thêm một lần nữa đã đúng chưa
	if !internal.PasswordStrong(password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": "Password must have at least 8 characters, " +
				"including uppercase letters, lowercase letters and numbers!",
		})
		return
	}

	// Băm mật khẩu
	hashedPassword, err := internal.HashPassword(password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	password = hashedPassword
	password = internal.Santize(password)
	email = internal.Santize(email)
	file, err := ctx.FormFile("file")
	if err != nil {
		newUser := &user_domain.User{
			ID:        primitive.NewObjectID(),
			IP:        ctx.ClientIP(),
			FullName:  fullName,
			AvatarURL: avatarUrl,
			Email:     email,
			Password:  hashedPassword,
			Verified:  false,
			Provider:  "fe-it",
			Role:      "user",
			Phone:     phone,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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

		var code string
		for {
			code = randstr.Dec(6)
			if u.UserUseCase.UniqueVerificationCode(ctx, code) {
				break
			}
		}

		updUser := user_domain.User{
			ID:               newUser.ID,
			VerificationCode: code,
			Verified:         false,
			UpdatedAt:        time.Now(),
		}

		// Update User in Database
		_, err = u.UserUseCase.UpdateVerify(ctx, &updUser)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error()},
			)
			return
		}

		emailData := google.EmailData{
			Code:      code,
			FirstName: newUser.FullName,
			Subject:   "Your account verification code",
		}

		err = google.SendEmail(&emailData, newUser.Email, "sign_in_first_time.html")
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "success",
				"message": "There was an error sending email",
			})
			return
		}

		// Trả về phản hồi thành công
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "We sent an email with a verification code to your email",
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

	imageURL, err := cloudinary.UploadImageToCloudinary(f, file.Filename, u.Database.CloudinaryUploadFolderUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	newUser := user_domain.User{
		ID:        primitive.NewObjectID(),
		FullName:  fullName,
		AvatarURL: imageURL.ImageURL,
		AssetID:   imageURL.AssetID,
		Email:     email,
		Password:  hashedPassword,
		Verified:  false,
		Provider:  "fe-it",
		Role:      "user",
		Phone:     phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// thực hiện đăng ký người dùng
	err = u.UserUseCase.Create(ctx, &newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	// add cron job
	//err = google.Cron.AddFunc("0h0m1s", func() {
	//	err = google.SendEmail(user.Email, subject_const.SignInTheFirstTime, subject_const.ContentTitle2)
	//	if err != nil {
	//		return
	//	}
	//})
	//google.Cron.Start()

	//if err != nil {
	//	return
	//}

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
	return
}

func (u *UserController) VerificationCode(ctx *gin.Context) {
	var verificationCode user_domain.VerificationCode
	if err := ctx.ShouldBindJSON(&verificationCode); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	user, err := u.UserUseCase.GetByVerificationCode(ctx, verificationCode.VerificationCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	res := u.UserUseCase.CheckVerify(ctx, verificationCode.VerificationCode)
	if res != true {
		ctx.JSON(http.StatusNotModified, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	updUser := user_domain.User{
		ID:        user.ID,
		Verified:  true,
		UpdatedAt: time.Now(),
	}

	// Update User in Database
	_, err = u.UserUseCase.UpdateVerify(ctx, &updUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error()},
		)
		return
	}

	// Trả về phản hồi thành công
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
