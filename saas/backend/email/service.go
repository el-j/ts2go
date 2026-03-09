package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/el-j/ts2go/saas/backend/logger"
)

// Service handles email sending
type Service struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// Config holds email service configuration
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// NewService creates a new email service
func NewService(cfg Config) *Service {
	return &Service{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		from:     cfg.From,
	}
}

// Email represents an email message
type Email struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

// Send sends an email
func (s *Service) Send(email *Email) error {
	// Build message
	headers := make(map[string]string)
	headers["From"] = s.from
	headers["To"] = email.To[0] // Simplified - should handle multiple recipients
	headers["Subject"] = email.Subject
	if email.IsHTML {
		headers["MIME-Version"] = "1.0"
		headers["Content-Type"] = "text/html; charset=UTF-8"
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + email.Body

	// Send email
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	err := smtp.SendMail(addr, auth, s.from, email.To, []byte(message))
	if err != nil {
		logger.Log.Error().
			Err(err).
			Strs("to", email.To).
			Str("subject", email.Subject).
			Msg("Failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.Log.Info().
		Strs("to", email.To).
		Str("subject", email.Subject).
		Msg("Email sent successfully")

	return nil
}

// SendJobComplete sends job completion notification
func (s *Service) SendJobComplete(to string, jobID string, projectName string, outputFiles int) error {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .button { display: inline-block; padding: 10px 20px; background: #4CAF50; color: white; text-decoration: none; border-radius: 5px; }
        .footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Transpilation Complete</h1>
        </div>
        <div class="content">
            <h2>Your TypeScript files have been successfully transpiled!</h2>
            <p><strong>Project:</strong> {{.ProjectName}}</p>
            <p><strong>Job ID:</strong> {{.JobID}}</p>
            <p><strong>Output Files:</strong> {{.OutputFiles}}</p>
            <p>Your transpiled Go files are ready for download.</p>
            <p style="text-align: center; margin-top: 30px;">
                <a href="https://ts2go.dev/projects/{{.ProjectName}}" class="button">View Results</a>
            </p>
        </div>
        <div class="footer">
            <p>TS2Go - TypeScript to Go Transpilation Service</p>
            <p>© 2025 TS2Go. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	data := struct {
		ProjectName string
		JobID       string
		OutputFiles int
	}{
		ProjectName: projectName,
		JobID:       jobID,
		OutputFiles: outputFiles,
	}

	var body bytes.Buffer
	t, err := template.New("jobComplete").Parse(tmpl)
	if err != nil {
		return err
	}

	if err := t.Execute(&body, data); err != nil {
		return err
	}

	return s.Send(&Email{
		To:      []string{to},
		Subject: fmt.Sprintf("Transpilation Complete - %s", projectName),
		Body:    body.String(),
		IsHTML:  true,
	})
}

// SendJobFailed sends job failure notification
func (s *Service) SendJobFailed(to string, jobID string, projectName string, errorMsg string) error {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #f44336; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .error { background: #ffebee; border-left: 4px solid #f44336; padding: 10px; margin: 20px 0; }
        .button { display: inline-block; padding: 10px 20px; background: #f44336; color: white; text-decoration: none; border-radius: 5px; }
        .footer { padding: 20px; text-align: center; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>❌ Transpilation Failed</h1>
        </div>
        <div class="content">
            <h2>We encountered an error processing your files</h2>
            <p><strong>Project:</strong> {{.ProjectName}}</p>
            <p><strong>Job ID:</strong> {{.JobID}}</p>
            <div class="error">
                <strong>Error:</strong><br>
                {{.ErrorMsg}}
            </div>
            <p>Please check your input files and try again. If the problem persists, contact support.</p>
            <p style="text-align: center; margin-top: 30px;">
                <a href="https://ts2go.dev/support" class="button">Contact Support</a>
            </p>
        </div>
        <div class="footer">
            <p>TS2Go - TypeScript to Go Transpilation Service</p>
            <p>© 2025 TS2Go. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	data := struct {
		ProjectName string
		JobID       string
		ErrorMsg    string
	}{
		ProjectName: projectName,
		JobID:       jobID,
		ErrorMsg:    errorMsg,
	}

	var body bytes.Buffer
	t, err := template.New("jobFailed").Parse(tmpl)
	if err != nil {
		return err
	}

	if err := t.Execute(&body, data); err != nil {
		return err
	}

	return s.Send(&Email{
		To:      []string{to},
		Subject: fmt.Sprintf("Transpilation Failed - %s", projectName),
		Body:    body.String(),
		IsHTML:  true,
	})
}

// SendWelcome sends welcome email to new users
func (s *Service) SendWelcome(to string, username string) error {
	body := fmt.Sprintf(`
Hello %s,

Welcome to TS2Go! We're excited to have you on board.

TS2Go is a powerful TypeScript to Go transpilation service that makes it easy to convert your TypeScript code to idiomatic Go.

Get started:
1. Create your first project
2. Upload your TypeScript files
3. Configure transpilation settings
4. Download your Go code

Visit https://ts2go.dev to get started!

Best regards,
The TS2Go Team
`, username)

	return s.Send(&Email{
		To:      []string{to},
		Subject: "Welcome to TS2Go!",
		Body:    body,
		IsHTML:  false,
	})
}
