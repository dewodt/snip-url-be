package emails

import (
	"os"

	"github.com/resend/resend-go/v2"
)

func generateWelcomeHtmlString() string {
	return `
		<!DOCTYPE html>
		<html>
		<head>
		</head>
		<body>
			<p>Welcome to Snip URL!</p>
		</body>
		</html>
	`
}

func SendWelcomeEmail(email string) (*resend.SendEmailResponse, error) {
	// Get resend secret
	apiKey := os.Getenv("RESEND_API_KEY")

	// Send email
	client := resend.NewClient(apiKey)
	params := &resend.SendEmailRequest{
		From:    "Snip URL <snip-url@dewodt.com>",
		To:      []string{email},
		Html:    generateWelcomeHtmlString(),
		Subject: "Welcome to Snip URL",
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		return nil, err
	}

	return sent, nil
}
