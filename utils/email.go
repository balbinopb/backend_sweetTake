package utils

import (
	"net/smtp"
	"os"
)

func SendResetEmail(toEmail, token string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte(
		"Subject: Sweetake Password Reset\n\n" +
			"Your password reset code:\n\n" +
			token + "\n\n" +
			"This code will expire in 15 minutes.",
	)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{toEmail},
		message,
	)
}
