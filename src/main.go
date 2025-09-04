package main

import (
	"disciplo/src/config"
	"disciplo/src/email"
	_ "disciplo/src/migrations"
	"disciplo/src/web"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

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
			if cfg.SMTPPort != "" {
				if cfg.SMTPPort == "587" {
					settings.SMTP.Port = 587
				} else if cfg.SMTPPort == "465" {
					settings.SMTP.Port = 465
				}
			}
			settings.SMTP.Username = cfg.SMTPUsername
			settings.SMTP.Password = cfg.SMTPPassword
			settings.SMTP.AuthMethod = "PLAIN"
			settings.SMTP.TLS = true
			
			if cfg.SMTPFrom != "" {
				settings.Meta.SenderName = cfg.SMTPFrom
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
				// Generate token and send admin invitation email
				token, _ := gonanoid.New(21)
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
		// Update admin user with telegram information
		admin, err := app.FindAuthRecordByEmail("users", cfg.AdminEmail)
		if err == nil && admin != nil {
			admin.Set("telegram_id", fmt.Sprintf("%d", message.From.ID))
			admin.Set("telegram_name", message.From.UserName)
			
			// Set verified to true when all telegram fields are filled
			telegramID := admin.GetString("telegram_id")
			telegramName := admin.GetString("telegram_name")
			if telegramID != "" && telegramName != "" {
				admin.Set("verified", true)
			}
			
			if err := app.Save(admin); err != nil {
				log.Printf("❌ Failed to update admin telegram info: %v", err)
			} else {
				log.Printf("✅ Admin telegram info updated in database")
			}
		}
		
		response = fmt.Sprintf("✅ **Admin Account Connected**\n\nToken: `%s`\n\n**Your Telegram Details:**\n• ID: `%d`\n• Username: @%s\n• Name: %s %s\n\nYour admin account is now linked to Telegram!",
			args, message.From.ID, message.From.UserName, message.From.FirstName, message.From.LastName)
		
		log.Printf("🔗 ADMIN LINKED - Token: %s | TG_ID: %d | Username: @%s | Name: %s %s",
			args, message.From.ID, message.From.UserName, message.From.FirstName, message.From.LastName)
	} else {
		response = "Welcome to **Disciplo**! 🎉\n\nTo connect your Telegram account, you need an invitation token from the admin.\n\n**How to get access:**\n1. Contact your administrator\n2. Get an invitation link\n3. Click the link to return here with a token"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = tgbotapi.ModeMarkdown

	// Add inline keyboard for dashboard access
	dashboardURL := cfg.Host + "/dashboard"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🌐 Visit Dashboard", dashboardURL),
		),
	)
	msg.ReplyMarkup = keyboard

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