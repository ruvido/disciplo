package main

import (
	"disciplo/src/config"
	"disciplo/src/email"
	_ "disciplo/src/migrations"
	"disciplo/src/web"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// parseSMTPFrom parses "Name <email@domain.com>" format
func parseSMTPFrom(smtpFrom string) (name, address string) {
	// Check if it matches "Name <email>" format
	if strings.Contains(smtpFrom, "<") && strings.Contains(smtpFrom, ">") {
		// Extract name (everything before <)
		parts := strings.SplitN(smtpFrom, "<", 2)
		name = strings.TrimSpace(strings.Trim(parts[0], "\""))
		
		// Extract address (between < and >)
		if len(parts) > 1 {
			address = strings.TrimSpace(strings.Trim(parts[1], ">"))
		}
	} else {
		// Just an email address without name
		name = ""
		address = strings.TrimSpace(smtpFrom)
	}
	return name, address
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := pocketbase.New()

	// Bootstrap hook for post-migration setup
	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		
		log.Printf("⚙️  Disciplo initialization complete")
		
		// Configure SMTP from environment variables
		if cfg.SMTPHost != "" && cfg.SMTPUsername != "" && cfg.SMTPPassword != "" {
			settings := e.App.Settings()
			settings.SMTP.Enabled = true
			settings.SMTP.Host = cfg.SMTPHost
			settings.SMTP.Port = 587
			settings.SMTP.Username = cfg.SMTPUsername
			settings.SMTP.Password = cfg.SMTPPassword
			settings.SMTP.AuthMethod = "PLAIN"
			settings.SMTP.TLS = false // Port 587 uses STARTTLS
			
			if cfg.SMTPFrom != "" {
				senderName, senderAddress := parseSMTPFrom(cfg.SMTPFrom)
				settings.Meta.SenderName = senderName
				settings.Meta.SenderAddress = senderAddress
			}
			
			if err := e.App.Save(settings); err != nil {
				log.Printf("⚠️  Failed to configure SMTP: %v", err)
			} else {
				log.Printf("✅ SMTP configured: %s:%s", cfg.SMTPHost, cfg.SMTPPort)
			}
		}
		
		// Check if admin was created by migration (check both collections)
		admin, err := e.App.FindAuthRecordByEmail("users", cfg.AdminEmail)
		superuser, _ := e.App.FindAuthRecordByEmail(core.CollectionNameSuperusers, cfg.AdminEmail)
		if err == nil && admin != nil && superuser != nil {
			log.Printf("✅ Admin %s ready (both users and superusers)", cfg.AdminEmail)
			
			// Only send email if admin is not yet verified (no telegram_id)
			if admin.GetString("telegram_id") == "" {
				// Generate and save admin token
				token, _ := gonanoid.New(21)
				admin.Set("telegram_token", token)
				if err := e.App.Save(admin); err != nil {
					log.Printf("⚠️  Failed to save admin telegram token: %v", err)
				}
				
				telegramLink := fmt.Sprintf("https://t.me/%s?start=%s", cfg.BotUsername, token)
				log.Printf("📱 Telegram Link: %s", telegramLink)
				
				// Send admin invitation email
				if err := email.SendAdminInvitation(e.App, cfg, telegramLink); err != nil {
					log.Printf("⚠️  Failed to send admin invitation email: %v", err)
				} else {
					log.Printf("📧 Admin invitation email sent to %s", cfg.AdminEmail)
				}
			} else {
				log.Printf("✅ Admin already verified with Telegram")
			}
		} else {
			if admin == nil {
				log.Printf("⚠️  Admin %s not found in users collection", cfg.AdminEmail)
			}
			if superuser == nil {
				log.Printf("⚠️  Admin %s not found in superusers collection", cfg.AdminEmail)
			}
		}
		
		return nil
	})

	// Configure port via OnServe hook
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Server.Addr = "0.0.0.0:" + cfg.Port
		return e.Next()
	})

	// Setup web routes
	web.SetupRoutes(app, cfg)

	// Start Telegram bot
	go startBot(app, cfg)

	log.Println("🚀 Starting Disciplo server...")
	log.Printf("🌐 Visit %s for dashboard", cfg.Host)
	log.Printf("📊 Admin panel: %s/_/", cfg.Host)
	
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}


func startBot(app core.App, cfg *config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Printf("❌ Bot failed to start: %v", err)
		log.Printf("💡 Check BOT_TOKEN in .env file")
		return
	}

	if cfg.DevMode {
		bot.Debug = true
	}

	log.Printf("✅ Telegram bot ready: @%s", bot.Self.UserName)

	// Update .env reminder if needed
	if cfg.BotUsername == "" || cfg.BotUsername == "your_bot_username" {
		log.Printf("💡 Update BOT_USERNAME in .env to: %s", bot.Self.UserName)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "start":
			handleStartCommand(bot, update.Message, app, cfg)
		case "help":
			handleHelpCommand(bot, update.Message)
		case "status":
			handleStatusCommand(bot, update.Message)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Use /help for available commands.")
			bot.Send(msg)
		}
	}
}

func handleStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, app core.App, cfg *config.Config) {
	args := message.CommandArguments()
	var response string
	
	if args != "" {
		// Find user by telegram_token (works for both admin and regular users)
		user, err := app.FindFirstRecordByFilter("users", "telegram_token = {:token}", map[string]interface{}{
			"token": args,
		})
		
		if err != nil || user == nil {
			response = "❌ **Invalid or Expired Token**\n\nThe token you used is not valid or has expired. Please contact your administrator for a new invitation link."
			log.Printf("❌ Invalid token used: %s", args)
		} else {
			
			// Update user with Telegram information
			user.Set("telegram_id", fmt.Sprintf("%d", message.From.ID))
			user.Set("telegram_name", message.From.UserName)
			user.Set("verified", true) // Now verified since Telegram is linked
			user.Set("telegram_token", "") // Clear the token after successful linking
			
			if err := app.Save(user); err != nil {
				log.Printf("❌ Failed to update user telegram info: %v", err)
				response = "❌ **Connection Failed**\n\nThere was an error linking your account. Please try again or contact support."
			} else {
				// Determine user role for message
				isAdmin := user.GetBool("admin")
				userName := user.GetString("name")
				
				// Use template for connection success message
				templateData := struct {
					FirstName string
					IsAdmin bool
				}{
					FirstName: message.From.FirstName,
					IsAdmin: isAdmin,
				}
				
				if templateMsg, err := loadBotTemplate("connection_success.md", templateData); err == nil {
					response = templateMsg
				} else {
					// Fallback message based on role
					if isAdmin {
						response = fmt.Sprintf("🎉 **Telegram Connected Successfully!**\n\nWelcome %s! Your Telegram account has been linked to Disciplo.\n\n✅ **Account Status**: Accepted\n👑 **Role**: Administrator\n🌐 **Community**: Disciplo\n\nYou can now access your dashboard to manage your profile and community settings.",
							message.From.FirstName)
					} else {
						response = fmt.Sprintf("🎉 **Telegram Connected Successfully!**\n\nWelcome %s! Your Telegram account has been linked to Disciplo.\n\n✅ **Account Status**: Verified\n👤 **Role**: Member\n🌐 **Community**: Disciplo\n\nYou are now verified and can access community features.",
							message.From.FirstName)
					}
				}
				
				log.Printf("🔗 USER LINKED - Name: %s | Email: %s | Admin: %v | TG_ID: %d | Username: @%s",
					userName, user.GetString("email"), isAdmin, message.From.ID, message.From.UserName)
			}
		}
	} else {
		response = "Welcome to **Disciplo**! 🎉\n\nTo connect your Telegram account, you need an invitation token from the admin.\n\n**How to get access:**\n1. Contact your administrator\n2. Get an invitation link\n3. Click the link to return here with a token"
	}

	// Add dashboard link logic before creating message
	dashboardURL := cfg.Host + "/dashboard"
	if !strings.HasPrefix(cfg.Host, "https://") {
		// For HTTP URLs, add dashboard link as text
		response += "\n\n🌐 **Dashboard**: " + dashboardURL
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = tgbotapi.ModeMarkdown

	// Add inline keyboard for dashboard access only if HTTPS
	if strings.HasPrefix(cfg.Host, "https://") {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("🌐 Visit Dashboard", dashboardURL),
			),
		)
		msg.ReplyMarkup = keyboard
	}

	bot.Send(msg)
}

func handleHelpCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	help := `**Disciplo Bot Commands**

🔗 **/start** - Connect account (requires invitation token)
❓ **/help** - Show this help message  
📊 **/status** - Check your account status

**Getting Started:**
Contact your community admin for an invitation link.`

	msg := tgbotapi.NewMessage(message.Chat.ID, help)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(msg)
}

func handleStatusCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	response := fmt.Sprintf("📊 **Your Account**\n\n"+
		"• **Telegram ID:** `%d`\n"+
		"• **Username:** @%s\n"+
		"• **Name:** %s %s\n\n"+
		"🔄 Account linking system is in development.\nContact admin for manual verification.",
		message.From.ID, message.From.UserName, message.From.FirstName, message.From.LastName)

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(msg)
}

// loadBotTemplate loads and renders a bot message template
func loadBotTemplate(templateName string, data interface{}) (string, error) {
	// Try to load template from file
	templatePath := filepath.Join("pb_public", "bot_templates", templateName)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}
	
	// Parse and execute template
	tmpl, err := template.New("bot").Parse(string(content))
	if err != nil {
		return "", err
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}