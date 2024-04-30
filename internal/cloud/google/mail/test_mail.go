package mail

import (
	"clean-architecture/internal/cloud/google"
	"github.com/thanhpk/randstr"
)

func SendMailTest() error {
	code := randstr.Dec(6)

	// ? Send Email
	emailData := google.EmailData{
		URL:     code,
		Subject: code + " Your account verification code",
	}

	//Thêm công việc cron để gửi email nhắc nhở
	err := google.SendEmail(&emailData, "verificationCode.html")
	if err != nil {
		return err
	}

	return nil
}
