package web

import (
	"disciplo/src/config"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
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

// Authentication helper function using PocketBase's native methods
func getAuthenticatedUser(c *core.RequestEvent) *core.Record {
	var token string
	
	// Check Authorization header first
	authorization := c.Request.Header.Get("Authorization")
	if authorization != "" {
		token = strings.TrimPrefix(authorization, "Bearer ")
	} else {
		// Try to get from cookie if no auth header  
		if cookie, err := c.Request.Cookie("pb_auth"); err == nil && cookie.Value != "" {
			token = cookie.Value
		}
	}
	
	if token == "" {
		return nil
	}
	
	// Use PocketBase's native FindAuthRecordByToken method
	// This automatically handles all auth collections
	record, err := c.App.FindAuthRecordByToken(token)
	if err != nil {
		return nil
	}
	
	return record
}

// Helper function for min operation
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Check if user is authenticated and is admin
func requireAdmin(c *core.RequestEvent) *core.Record {
	user := getAuthenticatedUser(c)
	if user == nil || !user.GetBool("admin") {
		return nil
	}
	return user
}

// Authentication middleware for regular users
func requireAuth(next func(*core.RequestEvent) error) func(*core.RequestEvent) error {
	return func(c *core.RequestEvent) error {
		user := getAuthenticatedUser(c)
		if user == nil {
			return c.Redirect(http.StatusFound, "/login")
		}
		return next(c)
	}
}

// Admin authentication middleware
func requireAdminAuth(next func(*core.RequestEvent) error) func(*core.RequestEvent) error {
	return func(c *core.RequestEvent) error {
		user := requireAdmin(c)
		if user == nil {
			return c.Redirect(http.StatusFound, "/login")
		}
		return next(c)
	}
}

type RegistrationData struct {
	PageTitle   string
	AppName     string
	Locations   []string
	JobFields   []string
	Interests   []string
}

type RegistrationRequest struct {
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Password       string   `json:"password"`
	DateOfBirth    string   `json:"date_of_birth"`
	City           string   `json:"city"`
	Location       string   `json:"location"`
	JobField       string   `json:"job_field"`
	Interests      []string `json:"interests"`
	WhyJoin        string   `json:"why_join"`
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
			// If user is already authenticated, redirect to appropriate dashboard
			if user := getAuthenticatedUser(c); user != nil {
				if user.GetBool("admin") {
					return c.Redirect(http.StatusFound, "/admin/dashboard")
				}
				return c.Redirect(http.StatusFound, "/dashboard")
			}

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

		// Logout route - clear session and redirect
		e.Router.GET("/logout", func(c *core.RequestEvent) error {
			// Clear authentication cookie
			cookie := &http.Cookie{
				Name:     "pb_auth",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				Expires:  time.Unix(0, 0),
				HttpOnly: false,
				SameSite: http.SameSiteStrictMode,
			}
			// Set secure flag only for HTTPS
			if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
				cookie.Secure = true
			}
			c.SetCookie(cookie)
			
			// Clear browser cache
			c.Response.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
			c.Response.Header().Set("Pragma", "no-cache")
			
			// Redirect to login page
			return c.Redirect(http.StatusFound, "/login")
		})
		
		// Logout API endpoint for AJAX requests
		e.Router.POST("/api/logout", func(c *core.RequestEvent) error {
			// Clear authentication cookie
			cookie := &http.Cookie{
				Name:     "pb_auth",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				Expires:  time.Unix(0, 0),
				HttpOnly: false,
				SameSite: http.SameSiteStrictMode,
			}
			// Set secure flag only for HTTPS
			if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
				cookie.Secure = true
			}
			c.SetCookie(cookie)
			
			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"message": "Logged out successfully",
			})
		})


		// User dashboard - PROTECTED
		e.Router.GET("/dashboard", func(c *core.RequestEvent) error {
			user := getAuthenticatedUser(c)
			if user == nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			userData := &UserData{
				Name:         user.GetString("name"),
				Email:        user.GetString("email"),
				Status:       user.GetString("status"),
				Admin:        user.GetBool("admin"),
				Verified:     user.GetBool("verified"),
				Created:      user.GetDateTime("created").String(),
				TelegramId:   user.GetString("telegram_id"),
				TelegramName: user.GetString("telegram_name"),
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

		// Profile page - PROTECTED
		e.Router.GET("/profile", func(c *core.RequestEvent) error {
			user := getAuthenticatedUser(c)
			if user == nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			userData := &UserData{
				Name:         user.GetString("name"),
				Email:        user.GetString("email"),
				Status:       user.GetString("status"),
				Admin:        user.GetBool("admin"),
				Verified:     user.GetBool("verified"),
				Created:      user.GetDateTime("created").String(),
				TelegramId:   user.GetString("telegram_id"),
				TelegramName: user.GetString("telegram_name"),
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

		// Communities page - PROTECTED (admin only)
		e.Router.GET("/communities", func(c *core.RequestEvent) error {
			user := requireAdmin(c)
			if user == nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			userData := &UserData{
				Name:         user.GetString("name"),
				Email:        user.GetString("email"),
				Status:       user.GetString("status"),
				Admin:        user.GetBool("admin"),
				Verified:     user.GetBool("verified"),
				Created:      user.GetDateTime("created").String(),
				TelegramId:   user.GetString("telegram_id"),
				TelegramName: user.GetString("telegram_name"),
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

		// Members page - PROTECTED (admin only)
		e.Router.GET("/members", func(c *core.RequestEvent) error {
			user := requireAdmin(c)
			if user == nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			userData := &UserData{
				Name:         user.GetString("name"),
				Email:        user.GetString("email"),
				Status:       user.GetString("status"),
				Admin:        user.GetBool("admin"),
				Verified:     user.GetBool("verified"),
				Created:      user.GetDateTime("created").String(),
				TelegramId:   user.GetString("telegram_id"),
				TelegramName: user.GetString("telegram_name"),
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

		// Admin dashboard - PROTECTED
		e.Router.GET("/admin/dashboard", func(c *core.RequestEvent) error {
			user := getAuthenticatedUser(c)
			if user == nil || !user.GetBool("admin") {
				return c.Redirect(http.StatusFound, "/login")
			}
			
			// Generate fresh token for admin
			token, _ := gonanoid.New(21)
			telegramLink := "https://t.me/" + cfg.BotUsername + "?start=" + token
			
			data := AdminData{
				AdminEmail:   user.GetString("email"),
				AdminName:    user.GetString("name"),
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
		
		// API endpoint to check telegram status - PROTECTED
		e.Router.GET("/api/admin/telegram-status", func(c *core.RequestEvent) error {
			user := getAuthenticatedUser(c)
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Not authenticated"})
			}
			
			telegramID := user.GetString("telegram_id")
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

		// API endpoint for profile updates - PROTECTED
		e.Router.PUT("/api/profile", func(c *core.RequestEvent) error {
			user := getAuthenticatedUser(c)
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Not authenticated",
				})
			}

			var updateData struct {
				Name string `json:"name"`
			}

			if err := c.BindBody(&updateData); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"error":   "Invalid request data",
				})
			}

			user.Set("name", updateData.Name)
			if err := e.App.Save(user); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"error":   "Failed to update profile",
				})
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		})

		// API endpoint for password change - PROTECTED
		e.Router.POST("/api/change-password", func(c *core.RequestEvent) error {
			user := getAuthenticatedUser(c)
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Not authenticated",
				})
			}

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

			// Verify current password
			if !user.ValidatePassword(passwordData.CurrentPassword) {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Current password is incorrect",
				})
			}

			// Set new password
			user.Set("password", passwordData.NewPassword)
			if err := e.App.Save(user); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"error":   "Failed to update password",
				})
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		})

		// Load disciplo configuration
		disciploConfig, err := config.LoadDisciploConfig()
		if err != nil {
			// Use default config if loading fails
			disciploConfig = &config.DisciploConfig{}
		}

		// Registration page
		e.Router.GET("/register", func(c *core.RequestEvent) error {
			if !disciploConfig.Registration.Enabled {
				return c.Redirect(http.StatusFound, "/login")
			}

			data := RegistrationData{
				PageTitle: "Register - " + disciploConfig.General.AppName,
				AppName:   disciploConfig.General.AppName,
				Locations: disciploConfig.Registration.Locations.Options,
				JobFields: disciploConfig.Registration.JobFields.Options,
				Interests: disciploConfig.Registration.Interests.Options,
			}

			tmpl, err := template.ParseFiles("pb_public/templates/register.html")
			if err != nil {
				return c.String(http.StatusInternalServerError, "Template error: "+err.Error())
			}

			var buf strings.Builder
			if err := tmpl.Execute(&buf, data); err != nil {
				return c.String(http.StatusInternalServerError, "Template error")
			}

			return c.HTML(http.StatusOK, buf.String())
		})

		// Registration API endpoint
		e.Router.POST("/api/register", func(c *core.RequestEvent) error {
			if !disciploConfig.Registration.Enabled {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"success": false,
					"error":   "Registration is currently disabled",
				})
			}

			// Parse form data
			name := c.Request.FormValue("name")
			email := c.Request.FormValue("email")
			password := c.Request.FormValue("password")
			dateOfBirth := c.Request.FormValue("date_of_birth")
			city := c.Request.FormValue("city")
			location := c.Request.FormValue("location")
			jobField := c.Request.FormValue("job_field")
			interestsJSON := c.Request.FormValue("interests")
			whyJoin := c.Request.FormValue("why_join")

			// Parse interests
			var interests []string
			if err := json.Unmarshal([]byte(interestsJSON), &interests); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"error":   "Invalid interests format",
				})
			}

			// Validate file upload exists (basic check)
			if c.Request.MultipartForm == nil || c.Request.MultipartForm.File["profile_picture"] == nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"error":   "Profile picture is required",
				})
			}

			// Basic file size validation
			fileHeaders := c.Request.MultipartForm.File["profile_picture"]
			if len(fileHeaders) > 0 && fileHeaders[0].Size > int64(disciploConfig.Registration.Picture.MaxSizeMB*1024*1024) {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"error":   fmt.Sprintf("File size must be less than %dMB", disciploConfig.Registration.Picture.MaxSizeMB),
				})
			}

			// Create new request record
			collection, err := e.App.FindCollectionByNameOrId("requests")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"error":   "Database error",
				})
			}

			record := core.NewRecord(collection)

			// Parse date of birth
			dob, err := time.Parse("2006-01-02", dateOfBirth)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"error":   "Invalid date of birth format",
				})
			}

			// Set record fields
			record.Set("name", name)
			record.Set("email", email)
			record.Set("password", password)
			dobDateTime := types.DateTime{}
			dobDateTime.Scan(dob)
			record.Set("date_of_birth", dobDateTime)
			record.Set("city", city)
			record.Set("location", location)
			record.Set("job_field", jobField)
			record.Set("interests", interests)
			record.Set("why_join", whyJoin)
			record.Set("status", "pending")

			// Save the record first to get an ID and allow file uploads
			if err := e.App.Save(record); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"error":   "Failed to save registration: " + err.Error(),
				})
			}

			// TODO: Process uploaded files using PocketBase
			// This will be implemented properly in the next iteration
			// For now, we'll skip file upload to get basic registration working

			// TODO: Send email notifications
			// This will be implemented in the next step

			return c.JSON(http.StatusOK, map[string]interface{}{
				"success":    true,
				"message":    "Registration submitted successfully",
				"request_id": record.Id,
			})
		})

		return e.Next()
	})
}