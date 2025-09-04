package main

import (
	"disciplo/src/bot"
	"disciplo/src/config"
	"disciplo/src/pb"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting Disciplo...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app, err := pb.InitializeSimple(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize PocketBase: %v", err)
	}

	if err := pb.CreateAdmin(app, cfg); err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	token := "test-token-123"
	if err := pb.SaveToken(app, cfg.AdminEmail, token); err != nil {
		log.Fatalf("Failed to save admin token: %v", err)
	}

	log.Printf("Admin setup complete. Token: %s", token)
	log.Printf("Telegram bot link: https://t.me/%s?start=%s", cfg.BotUsername, token)

	go bot.Start(app, cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		port := cfg.Port
		if port == "" {
			port = "8080"
		}
		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		}
		log.Printf("Starting server on http://localhost:%s", port)
		if err := app.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down...")
}