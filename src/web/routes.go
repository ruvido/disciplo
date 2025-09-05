package web

import (
	"disciplo/src/config"
	"disciplo/src/email"
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
	PageTitle      string
	AppName        string
	Locations      []string
	JobFields      []string
	Interests      []string
	RequiredFields map[string]bool
	Steps          []RegistrationStep
}

type RegistrationStep struct {
	Step   int      `json:"step"`
	Title  string   `json:"title"`
	Fields []string `json:"fields"`
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

// PocketBase middleware to redirect authenticated users from public pages
func redirectAuthenticatedUsers(next func(*core.RequestEvent) error) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		// Check if user is already authenticated using PocketBase patterns
		if user := getAuthenticatedUser(e); user != nil {
			// Redirect admin users to admin dashboard
			if user.GetBool("admin") {
				return e.Redirect(http.StatusFound, "/admin/dashboard")
			}
			// Redirect regular users to user dashboard
			return e.Redirect(http.StatusFound, "/dashboard")
		}
		// Continue to next handler if not authenticated
		return next(e)
	}
}

func SetupRoutes(app core.App, cfg *config.Config) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Root route - redirect authenticated users to dashboard, others to login
		e.Router.GET("/", redirectAuthenticatedUsers(func(c *core.RequestEvent) error {
			return c.Redirect(http.StatusFound, "/login")
		}))

		// Login page - use middleware to redirect authenticated users
		e.Router.GET("/login", redirectAuthenticatedUsers(func(c *core.RequestEvent) error {
			tmpl, err := template.ParseFiles("pb_public/templates/login.html")
			if err != nil {
				return c.String(http.StatusInternalServerError, "Template error: "+err.Error())
			}
			
			var buf strings.Builder
			if err := tmpl.Execute(&buf, nil); err != nil {
				return c.String(http.StatusInternalServerError, "Template error")
			}
			
			return c.HTML(http.StatusOK, buf.String())
		}))

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
			
			// Redirect admin users to admin dashboard
			if user.GetBool("admin") {
				return c.Redirect(http.StatusFound, "/admin/dashboard")
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

		// Admin requests page - PROTECTED
		e.Router.GET("/admin/requests", func(c *core.RequestEvent) error {
			user := getAuthenticatedUser(c)
			if user == nil || !user.GetBool("admin") {
				return c.Redirect(http.StatusFound, "/login")
			}

			// Load disciplo configuration
			disciploConfig, err := config.LoadDisciploConfig()
			if err != nil {
				disciploConfig = &config.DisciploConfig{}
			}

			// Get pending requests
			requests, err := e.App.FindRecordsByFilter("requests", "status = 'pending'", "", 50, 0)
			if err != nil {
				// Log error for debugging
				fmt.Printf("Error finding requests: %v\n", err)
				// Handle error but continue with empty list
				requests = []*core.Record{}
			} else {
				fmt.Printf("Found %d requests with status='pending'\n", len(requests))
			}

			data := struct {
				AdminEmail   string
				AdminName    string
				Requests     []*core.Record
				AppName      string
				Config       *config.DisciploConfig
			}{
				AdminEmail:   user.GetString("email"),
				AdminName:    user.GetString("name"),
				Requests:     requests,
				AppName:      disciploConfig.General.AppName,
				Config:       disciploConfig,
			}

			tmpl, err := template.ParseFiles("pb_public/templates/admin_requests.html")
			if err != nil {
				return c.String(http.StatusInternalServerError, "Template error: "+err.Error())
			}

			var buf strings.Builder
			if err := tmpl.Execute(&buf, data); err != nil {
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

		// API endpoint to approve membership request - ADMIN ONLY
		e.Router.POST("/api/admin/approve-request/{id}", func(c *core.RequestEvent) error {
			user := requireAdmin(c)
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Admin access required"})
			}

			requestId := c.Request.PathValue("id")
			if requestId == "" {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Request ID is required"})
			}
			
			// Find the request
			request, err := e.App.FindRecordById("requests", requestId)
			if err != nil {
				return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Request not found"})
			}

			// Check if request is already processed
			if request.GetString("status") != "pending" {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Request has already been processed"})
			}

			// Create user account in PocketBase's auth collection
			authCollection, err := e.App.FindCollectionByNameOrId("_pb_users_auth_")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to find auth collection"})
			}

			newUser := core.NewRecord(authCollection)
			
			// Set basic auth data from request
			newUser.Set("name", request.GetString("name"))
			newUser.Set("email", request.GetString("email"))
			newUser.Set("password", request.GetString("password"))
			newUser.Set("emailVisibility", false)
			// User starts with verified=false until Telegram is linked
			
			// Set additional fields for our platform
			newUser.Set("admin", false)
			newUser.Set("status", "accepted")
			newUser.Set("telegram_id", "")
			newUser.Set("telegram_name", "")
			newUser.Set("groups", "")
			newUser.Set("group_admin", "")
			newUser.Set("group_admin_since", "")

			// Save the new user
			if err := e.App.Save(newUser); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": "Failed to create user account: " + err.Error(),
				})
			}

			// Generate and save Telegram token for user
			token, _ := gonanoid.New(21)
			newUser.Set("telegram_token", token)
			if err := e.App.Save(newUser); err != nil {
				fmt.Printf("Warning: Failed to save telegram token: %v\n", err)
			}
			
			// Send welcome email with Telegram bot link
			if err := email.SendApprovalWelcome(e.App, request.GetString("email"), request.GetString("name"), cfg.BotUsername, token); err != nil {
				fmt.Printf("Warning: Failed to send welcome email: %v\n", err)
				// Continue with approval process even if email fails
			}

			// Update request with approval details
			request.Set("status", "approved")
			request.Set("approved_by", user.Id)
			request.Set("approved_at", types.NowDateTime())
			request.Set("created_user_id", newUser.Id)
			
			if err := e.App.Save(request); err != nil {
				// If request update fails, we should consider rolling back user creation
				// For simplicity, we'll log the error but continue
				fmt.Printf("Warning: Failed to update request after user creation: %v\n", err)
			}

			// TODO: Send welcome email to the new user
			// TODO: Send notification to admin about successful approval

			fmt.Printf("Request %s approved and user %s created successfully\n", requestId, newUser.Id)

			return c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"user_id": newUser.Id,
			})
		})

		// API endpoint to reject membership request - ADMIN ONLY
		e.Router.POST("/api/admin/reject-request/{id}", func(c *core.RequestEvent) error {
			user := requireAdmin(c)
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Admin access required"})
			}

			requestId := c.Request.PathValue("id")
			if requestId == "" {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Request ID is required"})
			}
			
			// Find the request
			request, err := e.App.FindRecordById("requests", requestId)
			if err != nil {
				return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Request not found"})
			}

			// Update request status to rejected
			request.Set("status", "rejected")
			if err := e.App.Save(request); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to update request"})
			}

			// TODO: Send rejection email (if implemented)

			return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
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

		// Registration page - use middleware to redirect authenticated users
		e.Router.GET("/register", redirectAuthenticatedUsers(func(c *core.RequestEvent) error {
			if !disciploConfig.Registration.Enabled {
				return c.Redirect(http.StatusFound, "/login")
			}

			// Build required fields map from configuration
			requiredFields := make(map[string]bool)
			var steps []RegistrationStep
			
			// Extract steps from disciplo config
			for _, step := range disciploConfig.Registration.Steps {
				steps = append(steps, RegistrationStep{
					Step:   step.Step,
					Title:  step.Title,
					Fields: step.Fields,
				})
				
				// Mark all fields in steps as required (as per config comment "all fields are required")
				for _, field := range step.Fields {
					requiredFields[field] = true
				}
			}
			
			// Profile picture is optional according to migration
			requiredFields["profile_picture"] = false

			data := RegistrationData{
				PageTitle:      "Register - " + disciploConfig.General.AppName,
				AppName:        disciploConfig.General.AppName,
				Locations:      disciploConfig.Registration.Locations.Options,
				JobFields:      disciploConfig.Registration.JobFields.Options,
				Interests:      disciploConfig.Registration.Interests.Options,
				RequiredFields: requiredFields,
				Steps:          steps,
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
		}))

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
			userEmail := c.Request.FormValue("email")
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

			// Optional file size validation (if file is uploaded)
			if c.Request.MultipartForm != nil && c.Request.MultipartForm.File["profile_picture"] != nil {
				fileHeaders := c.Request.MultipartForm.File["profile_picture"]
				if len(fileHeaders) > 0 && fileHeaders[0].Size > int64(disciploConfig.Registration.Picture.MaxSizeMB*1024*1024) {
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"success": false,
						"error":   fmt.Sprintf("File size must be less than %dMB", disciploConfig.Registration.Picture.MaxSizeMB),
					})
				}
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
			record.Set("email", userEmail)
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

			// Save the record  
			if err := e.App.Save(record); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"error":   "Failed to save registration: " + err.Error(),
				})
			}

			// Handle file upload - for now we skip complex file processing
			// Files will be handled by the form validation above
			// TODO: Implement proper file upload with PocketBase file handling

			// Send email notification to admin
			if err := email.SendNewRegistrationNotification(e.App, disciploConfig.General.EmailRequests, name, userEmail); err != nil {
				// Log error but don't fail the registration
				fmt.Printf("Failed to send admin notification email: %v\n", err)
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"success":    true,
				"message":    "Registration submitted successfully",
				"request_id": record.Id,
			})
		})

		// Email check API endpoint for duplicate validation
		e.Router.POST("/api/check-email", func(c *core.RequestEvent) error {
			email := c.Request.FormValue("email")
			if email == "" {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error": "Email is required",
				})
			}

			// Check if email exists in users collection
			userRecord, err := e.App.FindFirstRecordByFilter("users", "email = {:email}", map[string]interface{}{
				"email": email,
			})
			
			if err == nil && userRecord != nil {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"exists": true,
					"type":   "user",
					"message": "This email is already registered. Do you already have an account?",
				})
			}

			// Check if email exists in pending requests
			requestRecord, err := e.App.FindFirstRecordByFilter("requests", "email = {:email} AND status = 'pending'", map[string]interface{}{
				"email": email,
			})
			
			if err == nil && requestRecord != nil {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"exists": true,
					"type":   "pending",
					"message": "A registration request is already being processed for this email.",
				})
			}

			// Email is available
			return c.JSON(http.StatusOK, map[string]interface{}{
				"exists":  false,
				"type":    "none",
				"message": "",
			})
		})

		// Catch-all route for undefined paths - MUST BE LAST
		e.Router.GET("/*", func(c *core.RequestEvent) error {
			// Check authentication for undefined routes
			if user := getAuthenticatedUser(c); user != nil {
				// Authenticated user accessing undefined route → redirect to appropriate dashboard
				if user.GetBool("admin") {
					return c.Redirect(http.StatusFound, "/admin/dashboard")
				}
				return c.Redirect(http.StatusFound, "/dashboard")
			}
			// Unauthenticated user accessing undefined route → redirect to login
			return c.Redirect(http.StatusFound, "/login")
		})

		return e.Next()
	})
}