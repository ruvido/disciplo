package migrations

import (
	"fmt"
	"os"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		adminEmail := os.Getenv("ADMIN_EMAIL")
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		
		if adminEmail == "" || adminPassword == "" {
			return nil // Skip if env vars not set
		}

		// Create admin user in _superusers collection (for PocketBase admin access)
		existing, _ := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, adminEmail)
		if existing == nil {
			superusersCollection, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
			if err != nil {
				return err
			}
			
			superuser := core.NewRecord(superusersCollection)
			superuser.Set("email", adminEmail)
			superuser.Set("password", adminPassword)
			
			if err := app.Save(superuser); err != nil {
				return fmt.Errorf("failed to create superuser: %w", err)
			}
		}

		// Create admin user in users collection (for application access)
		usersCollection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		
		existingUser, _ := app.FindAuthRecordByEmail("users", adminEmail)
		if existingUser != nil {
			return nil // Admin already exists in users
		}

		// Create admin user in users collection
		record := core.NewRecord(usersCollection)
		record.Set("email", adminEmail)
		record.Set("password", adminPassword)
		record.Set("admin", true)
		record.Set("name", os.Getenv("ADMIN_NAME"))
		record.Set("status", "accepted")
		record.Set("verified", false)

		if err := app.Save(record); err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		return nil
	}, func(app core.App) error {
		// Revert operation
		adminEmail := os.Getenv("ADMIN_EMAIL")
		if adminEmail == "" {
			return nil
		}
		
		record, _ := app.FindAuthRecordByEmail("users", adminEmail)
		if record == nil {
			return nil
		}
		return app.Delete(record)
	})
}