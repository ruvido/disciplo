package email

import (
	"bytes"
	"disciplo/src/config"
	"fmt"
	"html/template"
	"net/mail"
	"os"
	"path/filepath"
	
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

// SendAdminInvitation sends welcome email to admin with Telegram link
func SendAdminInvitation(app core.App, cfg *config.Config, telegramLink string) error {
	// Load email template
	templatePath := filepath.Join("pb_public", "email_templates", "admin_invitation.html")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		// Fallback to default template
		content = []byte(getDefaultAdminInvitationTemplate())
	}
	
	// Parse template
	tmpl, err := template.New("email").Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}
	
	// Execute template with data
	var body bytes.Buffer
	data := map[string]interface{}{
		"Subject":     "Welcome to Disciplo - Connect Your Telegram",
		"AdminName":   cfg.AdminName,
		"AdminEmail":  cfg.AdminEmail,
		"BotUsername": cfg.BotUsername,
		"TelegramLink": telegramLink,
		"Host":        cfg.Host,
	}
	
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	
	// Send email using PocketBase mailer
	message := &mailer.Message{
		From: mail.Address{
			Address: app.Settings().Meta.SenderAddress,
			Name:    app.Settings().Meta.SenderName,
		},
		To: []mail.Address{{
			Address: cfg.AdminEmail,
		}},
		Subject: "Welcome to Disciplo - Connect Your Telegram",
		HTML:    body.String(),
	}
	
	return app.NewMailClient().Send(message)
}

func getDefaultAdminInvitationTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: -apple-system, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #f8f9fa; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: white; padding: 30px; border: 1px solid #e9ecef; }
        .button { display: inline-block; padding: 12px 24px; background: #0088cc; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { background: #f8f9fa; padding: 20px; text-align: center; font-size: 14px; color: #6c757d; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Disciplo!</h1>
        </div>
        <div class="content">
            <p>Hello {{.AdminName}},</p>
            <p>Your Disciplo platform is ready! To complete setup, please connect your Telegram account:</p>
            <p style="text-align: center;">
                <a href="{{.TelegramLink}}" class="button">Connect Telegram Account</a>
            </p>
            <p>Once connected, you can:</p>
            <ul>
                <li>Access the admin dashboard at {{.Host}}</li>
                <li>Manage community members</li>
                <li>Configure Telegram groups</li>
            </ul>
            <p>If the button doesn't work, copy this link: <br>{{.TelegramLink}}</p>
        </div>
        <div class="footer">
            <p>This is an automated message from Disciplo</p>
        </div>
    </div>
</body>
</html>`
}