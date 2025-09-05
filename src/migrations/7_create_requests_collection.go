package migrations

import (
	"disciplo/src/config"
	
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Check if requests collection already exists
		if collection, _ := app.FindCollectionByNameOrId("requests"); collection != nil {
			return nil // Collection exists, skip creation
		}

		// Load disciplo configuration
		disciploConfig, err := config.LoadDisciploConfig()
		if err != nil {
			// Use defaults if config loading fails - create basic collection
			disciploConfig = &config.DisciploConfig{}
		}

		// Create base collection
		collection := core.NewBaseCollection("requests")

		// Use location options from config
		locationOptions := []string{"Lazio", "Lombardia", "Piemonte", "Veneto", "Toscana"}
		if disciploConfig.Registration.Locations.Options != nil && len(disciploConfig.Registration.Locations.Options) > 0 {
			locationOptions = disciploConfig.Registration.Locations.Options
		}

		// Use job field options from config
		jobFieldOptions := []string{"Technology", "Finance", "Healthcare", "Education", "Other"}
		if disciploConfig.Registration.JobFields.Options != nil && len(disciploConfig.Registration.JobFields.Options) > 0 {
			jobFieldOptions = disciploConfig.Registration.JobFields.Options
		}

		// Use interests from config
		interestOptions := []string{"Networking", "Learning", "Technology", "Sport", "Cooking", "Reading"}
		if disciploConfig.Registration.Interests.Options != nil && len(disciploConfig.Registration.Interests.Options) > 0 {
			interestOptions = disciploConfig.Registration.Interests.Options
		}

		// Add all required fields
		collection.Fields = core.FieldsList{
			&core.TextField{
				Id:       "name",
				Name:     "name",
				Required: true,
			},
			&core.EmailField{
				Id:       "email",
				Name:     "email",
				Required: true,
			},
			&core.TextField{
				Id:       "password",
				Name:     "password",
				Required: true,
			},
			&core.DateField{
				Id:       "date_of_birth",
				Name:     "date_of_birth", 
				Required: true,
			},
			&core.TextField{
				Id:       "city",
				Name:     "city",
				Required: true,
			},
			&core.SelectField{
				Id:       "location",
				Name:     "location",
				Required: true,
				Values:   locationOptions,
			},
			&core.SelectField{
				Id:       "job_field",
				Name:     "job_field",
				Required: true,
				Values:   jobFieldOptions,
			},
			&core.SelectField{
				Id:       "interests",
				Name:     "interests",
				Required: true,
				Values:   interestOptions,
			},
			&core.TextField{
				Id:       "why_join",
				Name:     "why_join",
				Required: true,
			},
			&core.FileField{
				Id:       "profile_picture",
				Name:     "profile_picture",
				Required: false,
			},
			&core.SelectField{
				Id:       "status",
				Name:     "status",
				Required: true,
				Values:   []string{"pending", "approved", "rejected"},
			},
			&core.RelationField{
				Id:            "approved_by",
				Name:          "approved_by",
				CollectionId:  "_pb_users_auth_",
				CascadeDelete: false,
			},
			&core.DateField{
				Id:   "approved_at",
				Name: "approved_at",
			},
			&core.RelationField{
				Id:            "created_user_id",
				Name:          "created_user_id",
				CollectionId:  "_pb_users_auth_",
				CascadeDelete: false,
			},
			&core.AutodateField{
				Id:       "created",
				Name:     "created",
				OnCreate: true,
				OnUpdate: false,
			},
			&core.AutodateField{
				Id:       "updated",
				Name:     "updated",
				OnCreate: true,
				OnUpdate: true,
			},
		}

		return app.Save(collection)
	}, func(app core.App) error {
		// Revert - delete collection
		if collection, _ := app.FindCollectionByNameOrId("requests"); collection != nil {
			app.Delete(collection)
		}
		return nil
	})
}