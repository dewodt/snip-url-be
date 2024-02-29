package emails

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func generateSignInHtmlString(email string, token string) string {
	beUrl := os.Getenv("BE_URL")
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
		</head>
		<body>
			<p>Click the link below to sign in to Snip URL</p>
			<a href="%s/api/auth/email/callback?token=%s&email=%s">Sign in</a>
		</body>
		</html>	
	`, beUrl, token, email)
	return html
}

func SendSignInEmail(email string, token string) (*resend.SendEmailResponse, error) {
	// Get resend secret
	apiKey := os.Getenv("RESEND_API_KEY")

	// Send email
	client := resend.NewClient(apiKey)
	params := &resend.SendEmailRequest{
		From:    "Snip URL <snip-url@dewodt.com>",
		To:      []string{email},
		Html:    generateSignInHtmlString(email, token),
		Subject: "Sign in to Snip URL",
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		return nil, err
	}

	return sent, nil
}
