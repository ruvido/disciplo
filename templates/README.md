# Disciplo Templates System

This folder contains all templates for the Disciplo platform using Markdown with Go templating.

## Folder Structure

```
templates/
├── emails/          # Email templates (markdown + Go templates)
├── web/             # Future: Web page templates  
└── README.md        # This file
```

## Email Templates

All email templates are in `emails/` folder and use Markdown with Go template syntax.

### Available Templates

1. **new_request.md** - Sent to admin when new registration is submitted
2. **registration_received.md** - Sent to user when registration is submitted  
3. **approval_welcome.md** - Sent to user when registration is approved

### Template Variables

The following variables are available in all email templates:

#### Global Variables
- `{{.AppName}}` - Application name (from config)
- `{{.AdminEmail}}` - Admin email address
- `{{.BotUsername}}` - Telegram bot username
- `{{.Host}}` - Application host URL
- `{{.Date}}` - Current date/time

#### User Data Variables
- `{{.Name}}` - User's full name
- `{{.Email}}` - User's email address
- `{{.DateOfBirth}}` - User's date of birth
- `{{.City}}` - User's city
- `{{.Location}}` - User's location (region)
- `{{.JobField}}` - User's job field
- `{{.Interests}}` - Array of user's interests
- `{{.WhyJoin}}` - User's reason for joining
- `{{.ProfilePictureURL}}` - URL to user's profile picture

#### Request-Specific Variables
- `{{.RequestID}}` - Registration request ID
- `{{.Status}}` - Request status (pending/approved/rejected)
- `{{.SubmissionDate}}` - When request was submitted
- `{{.ApprovalDate}}` - When request was approved (if approved)

#### Admin-Specific Variables  
- `{{.ReviewURL}}` - Direct link to admin review page
- `{{.ApprovalURL}}` - Direct link to approve request
- `{{.DashboardURL}}` - Link to admin dashboard

#### Bot Integration Variables
- `{{.TelegramToken}}` - Token for Telegram bot connection
- `{{.BotDeepLink}}` - Complete Telegram bot deeplink with token

## Template Syntax

Templates use standard Go template syntax with Markdown formatting:

```markdown
# Welcome {{.Name}}!

Your registration has been **approved** for {{.AppName}}.

## Next Steps:
1. Connect your Telegram account: [Click Here]({{.BotDeepLink}})  
2. Access your dashboard: [Dashboard]({{.DashboardURL}})

## Your Application Details:
- **Email:** {{.Email}}
- **Location:** {{.City}}, {{.Location}}
- **Interests:** {{range .Interests}}{{.}}, {{end}}

Best regards,  
{{.AppName}} Team
```

## Adding New Variables

When adding new template variables:

1. Update the template engine in the Go code
2. Add the new variable to this README
3. Update existing templates if needed
4. Test with sample data

## Template Engine

Templates are processed by Go's `html/template` package with custom functions for:
- Markdown to HTML conversion
- Date formatting  
- Array/slice formatting
- URL generation

## Configuration

Email templates are configured in `disciplo.toml`:

```toml
[email]
template_engine = "markdown_go_template"
template_path = "templates/emails"

[email.templates]
new_request = "new_request.md"
registration_received = "registration_received.md"  
approval_welcome = "approval_welcome.md"
```