package main

import (
	"disciplo/src/config"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pocketbase/pocketbase"
)

func main() {
	log.Println("Starting Disciplo...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := pocketbase.New()

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if cfg.DevMode {
		bot.Debug = true
		log.Printf("Bot started as @%s", bot.Self.UserName)
	}

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		for update := range updates {
			if update.Message == nil {
				continue
			}

			if update.Message.IsCommand() && update.Message.Command() == "start" {
				response := fmt.Sprintf("✅ Welcome to %s!\n\nBot is running correctly.", cfg.AppName)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
				bot.Send(msg)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		port := cfg.Port
		if port == "" {
			port = "8080"
		}
		log.Printf("Starting server on http://localhost:%s", port)
		if err := app.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down...")
}