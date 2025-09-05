package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

// DisciploConfig represents the structure of disciplo.toml
type DisciploConfig struct {
	General      GeneralConfig      `toml:"general"`
	Registration RegistrationConfig `toml:"registration"`
	Email        EmailConfig        `toml:"email"`
	Admin        AdminConfig        `toml:"admin"`
	Auth         AuthConfig         `toml:"auth"`
}

type GeneralConfig struct {
	AppName       string `toml:"app_name"`
	EmailRequests string `toml:"email_requests"`
}

type RegistrationConfig struct {
	Enabled      bool           `toml:"enabled"`
	PublicAccess bool           `toml:"public_access"`
	Steps        []StepConfig   `toml:"steps"`
	JobFields    OptionsConfig  `toml:"job_fields"`
	Interests    OptionsConfig  `toml:"interests"`
	Locations    OptionsConfig  `toml:"locations"`
	Picture      PictureConfig  `toml:"picture"`
}

type StepConfig struct {
	Step   int      `toml:"step"`
	Title  string   `toml:"title"`
	Fields []string `toml:"fields"`
}

type OptionsConfig struct {
	Options []string `toml:"options"`
}

type PictureConfig struct {
	MaxSizeMB      int      `toml:"max_size_mb"`
	AllowedFormats []string `toml:"allowed_formats"`
	MaxDimensionPx int      `toml:"max_dimension_px"`
	AutoResize     bool     `toml:"auto_resize"`
}

type EmailConfig struct {
	TemplateEngine string           `toml:"template_engine"`
	TemplatePath   string           `toml:"template_path"`
	Templates      EmailTemplates   `toml:"templates"`
}

type EmailTemplates struct {
	NewRequest          string `toml:"new_request"`
	RegistrationReceived string `toml:"registration_received"`
	ApprovalWelcome     string `toml:"approval_welcome"`
}

type AdminConfig struct {
	RequestsPage        string `toml:"requests_page"`
	BulkOperations      bool   `toml:"bulk_operations"`
	IndividualReview    bool   `toml:"individual_review_only"`
	RejectionReasons    bool   `toml:"rejection_reasons"`
}

type AuthConfig struct {
	PendingUsersCanLogin       bool `toml:"pending_users_can_login"`
	ShowSignupLinkWhenNoUser   bool `toml:"show_signup_link_when_no_user"`
	ShowWaitMessageWhenPending bool `toml:"show_wait_message_when_pending"`
}

// LoadDisciploConfig loads configuration from disciplo.toml file
func LoadDisciploConfig() (*DisciploConfig, error) {
	var config DisciploConfig
	
	// Try to read disciplo.toml from current directory
	tomlData, err := os.ReadFile("disciplo.toml")
	if err != nil {
		// Try build directory
		tomlData, err = os.ReadFile("build/disciplo.toml")
		if err != nil {
			// Return default config if file not found
			return getDefaultConfig(), nil
		}
	}

	if err := toml.Unmarshal(tomlData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getDefaultConfig() *DisciploConfig {
	return &DisciploConfig{
		General: GeneralConfig{
			AppName:       "Disciplo",
			EmailRequests: "admin@example.com",
		},
		Registration: RegistrationConfig{
			Enabled:      true,
			PublicAccess: true,
			JobFields: OptionsConfig{
				Options: []string{"Technology", "Finance", "Healthcare", "Education", "Marketing", "Consulting", "Other"},
			},
			Interests: OptionsConfig{
				Options: []string{"Networking", "Learning", "Mentoring", "Collaboration"},
			},
			Locations: OptionsConfig{
				Options: []string{"Lazio", "Lombardia", "Piemonte", "Veneto", "Toscana"},
			},
			Picture: PictureConfig{
				MaxSizeMB:      5,
				AllowedFormats: []string{"jpg", "jpeg", "png"},
				MaxDimensionPx: 400,
				AutoResize:     true,
			},
		},
		Email: EmailConfig{
			TemplateEngine: "markdown_go_template",
			TemplatePath:   "templates/emails",
			Templates: EmailTemplates{
				NewRequest:          "new_request.md",
				RegistrationReceived: "registration_received.md",
				ApprovalWelcome:     "approval_welcome.md",
			},
		},
		Admin: AdminConfig{
			RequestsPage:        "/admin/requests",
			BulkOperations:      false,
			IndividualReview:    true,
			RejectionReasons:    false,
		},
		Auth: AuthConfig{
			PendingUsersCanLogin:       false,
			ShowSignupLinkWhenNoUser:   true,
			ShowWaitMessageWhenPending: true,
		},
	}
}