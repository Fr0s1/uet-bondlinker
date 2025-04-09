
package util

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
	"socialnet/config"
)

// EmailData holds data to be used in email templates
type EmailData struct {
	Name          string
	Email         string
	Subject       string
	VerifyLink    string
	ResetLink     string
	SiteBaseURL   string
	ExpiryMinutes int
}

// EmailService provides email functionality
type EmailService struct {
	cfg *config.Config
}

// NewEmailService creates a new EmailService
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

// SendWelcomeEmail sends a welcome email to a new user
func (es *EmailService) SendWelcomeEmail(name, email, token string) error {
	data := EmailData{
		Name:        name,
		Email:       email,
		Subject:     "Welcome to SocialNet - Verify Your Email",
		VerifyLink:  fmt.Sprintf("%s/verify-email?token=%s", es.cfg.Email.FrontendURL, token),
		SiteBaseURL: es.cfg.Email.FrontendURL,
	}

	return es.SendEmail(email, "Welcome to SocialNet - Verify Your Email", "welcome", data)
}

// SendPasswordResetEmail sends an email with password reset instructions
func (es *EmailService) SendPasswordResetEmail(name, email, token string) error {
	data := EmailData{
		Name:          name,
		Email:         email,
		Subject:       "SocialNet - Reset Your Password",
		ResetLink:     fmt.Sprintf("%s/reset-password?token=%s", es.cfg.Email.FrontendURL, token),
		SiteBaseURL:   es.cfg.Email.FrontendURL,
		ExpiryMinutes: 15,
	}

	return es.SendEmail(email, "SocialNet - Reset Your Password", "reset-password", data)
}

// SendEmail sends an email with the specified template
func (es *EmailService) SendEmail(to, subject, templateName string, data EmailData) error {
	from := es.cfg.Email.FromEmail
	password := es.cfg.Email.Password
	host := es.cfg.Email.SMTPHost
	port := es.cfg.Email.SMTPPort
	address := fmt.Sprintf("%s:%s", host, port)

	// Create authentication
	auth := smtp.PlainAuth("", from, password, host)

	// Get email template
	templatePath := filepath.Join("templates", "emails", templateName+".html")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		// If template file not found, use fallback template from string
		fallbackTemplate := getFallbackTemplate(templateName)
		t, err = template.New(templateName).Parse(fallbackTemplate)
		if err != nil {
			return err
		}
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	// Create email headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Construct email message
	var message bytes.Buffer
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.Write(body.Bytes())

	// Send email
	return smtp.SendMail(address, auth, from, []string{to}, message.Bytes())
}

// getFallbackTemplate returns a basic email template if file not found
func getFallbackTemplate(templateType string) string {
	switch templateType {
	case "welcome":
		return `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Welcome to SocialNet</title>
		</head>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
			<div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin-bottom: 20px;">
				<h1 style="color: #3b82f6;">Welcome to SocialNet!</h1>
			</div>
			<p>Hi {{.Name}},</p>
			<p>Thank you for joining SocialNet. To complete your registration, please verify your email address by clicking the button below:</p>
			<div style="text-align: center; margin: 30px 0;">
				<a href="{{.VerifyLink}}" style="background-color: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; font-weight: bold;">Verify My Email</a>
			</div>
			<p>If the button doesn't work, you can also copy and paste this link into your browser:</p>
			<p><a href="{{.VerifyLink}}">{{.VerifyLink}}</a></p>
			<p>If you didn't sign up for SocialNet, you can safely ignore this email.</p>
			<hr style="margin: 30px 0; border: none; border-top: 1px solid #eee;">
			<p style="color: #666; font-size: 14px;">The SocialNet Team</p>
		</body>
		</html>
		`
	case "reset-password":
		return `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Reset Your Password</title>
		</head>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
			<div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin-bottom: 20px;">
				<h1 style="color: #3b82f6;">Reset Your Password</h1>
			</div>
			<p>Hi {{.Name}},</p>
			<p>We received a request to reset your password. Click the button below to create a new password:</p>
			<div style="text-align: center; margin: 30px 0;">
				<a href="{{.ResetLink}}" style="background-color: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; font-weight: bold;">Reset Password</a>
			</div>
			<p>If the button doesn't work, you can also copy and paste this link into your browser:</p>
			<p><a href="{{.ResetLink}}">{{.ResetLink}}</a></p>
			<p>This link will expire in {{.ExpiryMinutes}} minutes.</p>
			<p>If you didn't request a password reset, you can safely ignore this email.</p>
			<hr style="margin: 30px 0; border: none; border-top: 1px solid #eee;">
			<p style="color: #666; font-size: 14px;">The SocialNet Team</p>
		</body>
		</html>
		`
	default:
		return `
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{.Subject}}</title>
		</head>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
			<div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin-bottom: 20px;">
				<h1 style="color: #3b82f6;">{{.Subject}}</h1>
			</div>
			<p>Hi {{.Name}},</p>
			<p>This is an automated email from SocialNet.</p>
			<hr style="margin: 30px 0; border: none; border-top: 1px solid #eee;">
			<p style="color: #666; font-size: 14px;">The SocialNet Team</p>
		</body>
		</html>
		`
	}
}

// Create email templates directory if it doesn't exist
func EnsureEmailTemplatesDir() error {
	templatesDir := filepath.Join("templates", "emails")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return err
	}

	// Create sample templates if they don't exist
	templates := map[string]string{
		"welcome.html":       getFallbackTemplate("welcome"),
		"reset-password.html": getFallbackTemplate("reset-password"),
	}

	for filename, content := range templates {
		filePath := filepath.Join(templatesDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
