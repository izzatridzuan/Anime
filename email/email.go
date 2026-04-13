package email

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendWelcomeEmail(toEmail, toName, password string) error {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	appURL := os.Getenv("APP_URL")

	body := fmt.Sprintf(`
        <h2>Welcome to Anime Backoffice, %s!</h2>
        <p>Your account has been created. Here are your login details:</p>
        <p><strong>Email:</strong> %s</p>
        <p><strong>Password:</strong> %s</p>
        <p>You will be prompted to change your password on first login.</p>
        <br/>
        <p>Login at: <a href="%s/login">%s/login</a></p>
    `, toName, toEmail, password, appURL, appURL)

	params := &resend.SendEmailRequest{
		From:    os.Getenv("RESEND_FROM"),
		To:      []string{toEmail},
		Subject: "Your Anime Backoffice Account",
		Html:    body,
	}

	_, err := client.Emails.Send(params)
	return err
}
