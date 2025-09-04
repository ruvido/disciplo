package bot

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// LoadTemplate loads and parses a bot template file
func LoadTemplate(name string, data interface{}) (string, error) {
	// Try to load from pb_public first (runtime), then from source
	paths := []string{
		filepath.Join("pb_public", "bot_templates", name),
		filepath.Join("src", "static", "bot_templates", name),
	}
	
	var content []byte
	var err error
	
	for _, path := range paths {
		content, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}
	
	if err != nil {
		return "", fmt.Errorf("template %s not found", name)
	}
	
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	
	return buf.String(), nil
}

// GetDefaultTemplate returns a hardcoded fallback template
func GetDefaultTemplate(name string) string {
	templates := map[string]string{
		"start_admin.md": `✅ **Admin Account Connected**

Token: {{.Token}}

**Your Telegram Details:**
• ID: {{.TelegramID}}
• Username: @{{.Username}}
• Name: {{.FirstName}} {{.LastName}}

Your admin account is now linked to Telegram!`,
		"start_welcome.md": `Welcome to **Disciplo**! 🎉

To connect your Telegram account, you need an invitation token from the admin.

**How to get access:**
1. Contact your administrator
2. Get an invitation link
3. Click the link to return here with a token`,
		"help.md": `**Disciplo Bot Commands**

🔗 **/start** - Connect account (requires invitation token)
❓ **/help** - Show this help message  
📊 **/status** - Check your account status

**Getting Started:**
Contact your community admin for an invitation link.`,
		"status.md": `📊 **Your Account**

• **Telegram ID:** {{.TelegramID}}
• **Username:** @{{.Username}}
• **Name:** {{.FirstName}} {{.LastName}}

🔄 Account linking system is in development.
Contact admin for manual verification.`,
	}
	
	if tmpl, ok := templates[name]; ok {
		return tmpl
	}
	return ""
}