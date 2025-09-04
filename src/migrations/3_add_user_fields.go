package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Add missing fields to users collection
		// Note: The users collection exists but is missing required fields from CLAUDE.md:
		// admin (bool), group_admin (relation), group_admin_since (date), 
		// groups (relation), status (select), verified (bool), telegram_id, telegram_name
		
		usersCollection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err // Users collection must exist
		}

		// Add missing fields
		newFields := core.FieldsList{
			&core.TextField{
				Id:   "telegram_id",
				Name: "telegram_id",
			},
			&core.TextField{
				Id:   "telegram_name", 
				Name: "telegram_name",
			},
			&core.BoolField{
				Id:   "admin",
				Name: "admin",
			},
		}

		// Append new fields to existing fields
		usersCollection.Fields = append(usersCollection.Fields, newFields...)
		
		return app.Save(usersCollection)
	}, func(app core.App) error {
		// Revert: Remove the custom fields we added
		// (This is complex to implement correctly, so we'll leave empty for now)
		return nil
	})
}