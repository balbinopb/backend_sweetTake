package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

// func SendResetEmail(toEmail, token string) error {
// 	from := os.Getenv("SMTP_EMAIL")
// 	password := os.Getenv("SMTP_PASSWORD")

// 	smtpHost := os.Getenv("SMTP_HOST")
// 	smtpPort := os.Getenv("SMTP_PORT")

// 	auth := smtp.PlainAuth("", from, password, smtpHost)

// 	message := []byte(
// 		"Subject: Sweetake Password Reset\n\n" +
// 			"Your password reset code:\n\n" +
// 			token + "\n\n" +
// 			"This code will expire in 15 minutes.",
// 	)

//		return smtp.SendMail(
//			smtpHost+":"+smtpPort,
//			auth,
//			from,
//			[]string{toEmail},
//			message,
//		)
//	}
func SendResetEmail(email, token string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return errors.New("RESEND_API_KEY not set")
	}

	body := map[string]interface{}{
		"from":    "SweetTake <onboarding@resend.dev>",
		"to":      []string{email},
		"subject": "Password Reset",
		"html":    "<p>Your reset token is:</p><h2>" + token + "</h2><p>Valid for 15 minutes.</p>",
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"POST",
		"https://api.resend.com/emails",
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("failed to send email")
	}

	return nil
}
