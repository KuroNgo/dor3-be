package google

import (
	subject_const "clean-architecture/internal/cloud/google/const"
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader(subject_const.From, subject_const.Mailer1, subject_const.Password1)
	m.SetHeader(subject_const.To, to)

	m.SetAddressHeader(subject_const.Bcc, subject_const.BCCAdmin1, subject_const.Admin)
	m.SetAddressHeader(subject_const.Bcc, subject_const.BCCAdmin2, subject_const.Admin)
	m.SetAddressHeader(subject_const.Bcc, subject_const.BCCAdmin3, subject_const.Admin)

	m.SetHeader(subject_const.Subject, subject)
	m.SetBody(
		subject_const.Body,
		"<h3 style=\"font-family: 'Marker Felt', sans-serif; font-size: 20px; line-height: 1.5; color: #00DD00;\"><strong>"+body+"</strong></h3>",
	)

	// random image
	//m.Attach("/C:/Users/hoaip/Downloads/wallpaperflare.com_wallpaper.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, subject_const.Mailer1, subject_const.Password1)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
