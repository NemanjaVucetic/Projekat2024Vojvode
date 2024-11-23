package utils

import (
	"fmt"
	"userService/domain"

	"gopkg.in/gomail.v2"
)

func SendEmail(user *domain.User, subject string) error {
	validationLink := fmt.Sprintf("http://localhost:8000/api/user/validate/%s", user.ID.Hex())

	body := fmt.Sprintf(`
		<p>Dear %s,</p>
		<p>Thank you for registering. Please confirm your account by clicking the link below:</p>
		<a href="%s">Confirm Account</a>
		`, user.Name, validationLink)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "mailexample81@gmail.com")
	mailer.SetHeader("To", user.Email)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "mailexample81@gmail.com", "nmpm qsgm dixv xkqg") // Use your App Password if 2FA enabled

	return dialer.DialAndSend(mailer)
}
