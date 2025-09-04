package web

import (
	"disciplo/src/config"
	"html/template"
	"net/http"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pocketbase/pocketbase/core"
)

type AdminData struct {
	AdminEmail   string
	AdminName    string
	TelegramLink string
}

type UserData struct {
	Name         string
	Email        string
	Status       string
	Admin        bool
	Verified     bool
	Created      string
	TelegramId   string
	TelegramName string
}

type PageData struct {
	PageTitle   string
	AppName     string
	User        *UserData
	BotUsername string
}

func SetupRoutes(app core.App, cfg *config.Config) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Root route - redirect to login if not authenticated
		e.Router.GET("/", func(c *core.RequestEvent) error {
			// For now, always redirect to login - auth will be handled by client-side
			return c.Redirect(http.StatusFound, "/login")
		})

		// Login page
		e.Router.GET("/login", func(c *core.RequestEvent) error {
			tmpl, err := template.ParseFiles("pb_public/templates/login.html")
			if err != nil {
				return c.String(http.StatusInternalServerError, "Template error: "+err.Error())
			}
			
			var buf strings.Builder
			if err := tmpl.Execute(&buf, nil); err != nil {
				return c.String(http.StatusInternalServerError, "Template error")
			}
			
			return c.HTML(http.StatusOK, buf.String())
		})

		// Logout route
		e.Router.GET("/logout", func(c *core.RequestEvent) error {
			// Clear any authentication cookies/session
			cookie := &http.Cookie{
				Name:     "pb_auth",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			}
			c.SetCookie(cookie)
			
			// Redirect to login page
			return c.Redirect(http.StatusFound, "/login")
		})


		// User dashboard
		e.Router.GET("/dashboard", func(c *core.RequestEvent) error {
			// Get user data (for now using admin as example - would need proper auth in real app)
			admin, err := e.App.FindAuthRecordByEmail("users", cfg.AdminEmail)
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			userData := &UserData{
				Name:         admin.GetString("name"),
				Email:        admin.GetString("email"),
				Status:       admin.GetString("status"),
				Admin:        admin.GetBool("admin"),
				Verified:     admin.GetBool("verified"),
				Created:      admin.GetDateTime("created").String(),
				TelegramId:   admin.GetString("telegram_id"),
				TelegramName: admin.GetString("telegram_name"),
			}

			pageData := PageData{
				PageTitle:   "Dashboard",
				AppName:     cfg.AppName,
				User:        userData,
				BotUsername: cfg.BotUsername,
			}

			tmpl, err := template.ParseFiles("pb_public/templates/base.html", "pb_public/templates/dashboard.html")
			if err != nil {
				return c.String(http.StatusInternalServerError, "Template error: "+err.Error())
			}
			
			var buf strings.Builder
			if err := tmpl.Execute(&buf, pageData); err != nil {
				return c.String(http.StatusInternalServerError, "Template error")
			}
			
			return c.HTML(http.StatusOK, buf.String())
		})

		// Profile page
		e.Router.GET("/profile", func(c *core.RequestEvent) error {
			// Get user data (for now using admin as example - would need proper auth in real app)
			admin, err := e.App.FindAuthRecordByEmail("users", cfg.AdminEmail)
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			userData := &UserData{
				Name:         admin.GetString("name"),
				Email:        admin.GetString("email"),
				Status:       admin.GetString("status"),
				Admin:        admin.GetBool("admin"),
				Verified:     admin.GetBool("verified"),
				Created:      admin.GetDateTime("created").String(),
				TelegramId:   admin.GetString("telegram_id"),
				TelegramName: admin.GetString("telegram_name"),
			}

			pageData := PageData{
				PageTitle:   "Profile",
				AppName:     cfg.AppName,
				User:        userData,
				BotUsername: cfg.BotUsername,
			}

			tmpl, err := template.ParseFiles("pb_public/templates/base.html", "pb_public/templates/profile.html")
			if err != nil {
				return c.String(http.StatusInternalServerError, "Template error: "+err.Error())
			}
			
			var buf strings.Builder
			if err := tmpl.Execute(&buf, pageData); err != nil {
				return c.String(http.StatusInternalServerError, "Template error")
			}
			
			return c.HTML(http.StatusOK, buf.String())
		})

		// Communities page
		e.Router.GET("/communities", func(c *core.RequestEvent) error {
			// Get user data (for now using admin as example - would need proper auth in real app)
			admin, err := e.App.FindAuthRecordByEmail("users", cfg.AdminEmail)
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			userData := &UserData{
				Name:         admin.GetString("name"),
				Email:        admin.GetString("email"),
				Status:       admin.GetString("status"),
				Admin:        admin.GetBool("admin"),
				Verified:     admin.GetBool("verified"),
				Created:      admin.GetDateTime("created").String(),
				TelegramId:   admin.GetString("telegram_id"),
				TelegramName: admin.GetString("telegram_name"),
			}

			pageData := PageData{
				PageTitle:   "Communities",
				AppName:     cfg.AppName,
				User:        userData,
				BotUsername: cfg.BotUsername,
			}

			tmpl, err := template.ParseFiles("pb_public/templates/base.html", "pb_public/templates/communities.html")
			if err != nil {
				return c.String(http.StatusInternalServerError, "Template error: "+err.Error())
			}
			
			var buf strings.Builder
			if err := tmpl.Execute(&buf, pageData); err != nil {
				return c.String(http.StatusInternalServerError, "Template error")
			}
			
			return c.HTML(http.StatusOK, buf.String())
		})

		// Members page
		e.Router.GET("/members", func(c *core.RequestEvent) error {
			// Get user data (for now using admin as example - would need proper auth in real app)
			admin, err := e.App.FindAuthRecordByEmail("users", cfg.AdminEmail)
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			userData := &UserData{
				Name:         admin.GetString("name"),
				Email:        admin.GetString("email"),
				Status:       admin.GetString("status"),
				Admin:        admin.GetBool("admin"),
				Verified:     admin.GetBool("verified"),
				Created:      admin.GetDateTime("created").String(),
				TelegramId:   admin.GetString("telegram_id"),
				TelegramName: admin.GetString("telegram_name"),
			}

			pageData := PageData{
				PageTitle:   "Members",
				AppName:     cfg.AppName,
				User:        userData,
				BotUsername: cfg.BotUsername,
			}

			tmpl, err := template.ParseFiles("pb_public/templates/base.html", "pb_public/templates/members.html")
			if err != nil {
				return c.String(http.StatusInternalServerError, "Template error: "+err.Error())
			}
			
			var buf strings.Builder
			if err := tmpl.Execute(&buf, pageData); err != nil {
				return c.String(http.StatusInternalServerError, "Template error")
			}
			
			return c.HTML(http.StatusOK, buf.String())
		})

		// Admin dashboard
		e.Router.GET("/admin/dashboard", func(c *core.RequestEvent) error {
			
			// Generate fresh token for admin
			token, _ := gonanoid.New(21)
			telegramLink := "https://t.me/" + cfg.BotUsername + "?start=" + token
			
			data := AdminData{
				AdminEmail:   cfg.AdminEmail,
				AdminName:    cfg.AdminName,
				TelegramLink: telegramLink,
			}

			tmpl, err := template.ParseFiles("pb_public/templates/admin_dashboard.html")
			if err != nil {
				return c.String(http.StatusInternalServerError, "Template error: "+err.Error())
			}

			var buf strings.Builder
			if err := tmpl.Execute(&buf, data); err != nil {
				return c.String(http.StatusInternalServerError, "Template error")
			}

			return c.HTML(http.StatusOK, buf.String())
		})
		
		// API endpoint to check telegram status
		e.Router.GET("/api/admin/telegram-status", func(c *core.RequestEvent) error {
			admin, err := e.App.FindAuthRecordByEmail("users", cfg.AdminEmail)
			if err != nil {
				return c.JSON(http.StatusOK, map[string]bool{"connected": false})
			}
			
			telegramID := admin.GetString("telegram_id")
			connected := telegramID != ""
			
			return c.JSON(http.StatusOK, map[string]bool{"connected": connected})
		})

		// API endpoint to generate token for telegram connection
		e.Router.POST("/api/generate-token", func(c *core.RequestEvent) error {
			// For now, generate a simple token - in a real app you'd want to store this temporarily
			token, err := gonanoid.New(21)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"error":   "Failed to generate token",
				})
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"token":   token,
			})
		})

		// API endpoint for profile updates
		e.Router.PUT("/api/profile", func(c *core.RequestEvent) error {
			// This would need proper authentication in a real app
			// For now, we'll assume this is for the admin user
			var updateData struct {
				Name string `json:"name"`
			}

			if err := c.BindBody(&updateData); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"error":   "Invalid request data",
				})
			}

			admin, err := e.App.FindAuthRecordByEmail("users", cfg.AdminEmail)
			if err != nil {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"success": false,
					"error":   "User not found",
				})
			}

			admin.Set("name", updateData.Name)
			if err := e.App.Save(admin); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"error":   "Failed to update profile",
				})
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		})

		// API endpoint for password change
		e.Router.POST("/api/change-password", func(c *core.RequestEvent) error {
			var passwordData struct {
				CurrentPassword string `json:"currentPassword"`
				NewPassword     string `json:"newPassword"`
			}

			if err := c.BindBody(&passwordData); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"error":   "Invalid request data",
				})
			}

			if len(passwordData.NewPassword) < 8 {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"error":   "New password must be at least 8 characters",
				})
			}

			admin, err := e.App.FindAuthRecordByEmail("users", cfg.AdminEmail)
			if err != nil {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"success": false,
					"error":   "User not found",
				})
			}

			// Verify current password
			if !admin.ValidatePassword(passwordData.CurrentPassword) {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Current password is incorrect",
				})
			}

			// Set new password
			admin.Set("password", passwordData.NewPassword)
			if err := e.App.Save(admin); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"error":   "Failed to update password",
				})
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		})

		return e.Next()
	})
}