package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Add telegram_token field to users collection
		usersCollection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err // Users collection must exist
		}

		// Add telegram_token field
		telegramTokenField := &core.TextField{
			Id:   "telegram_token",
			Name: "telegram_token",
		}

		// Append new field to existing fields
		usersCollection.Fields = append(usersCollection.Fields, telegramTokenField)
		
		return app.Save(usersCollection)
	}, func(app core.App) error {
		// Revert: Remove the telegram_token field
		usersCollection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		
		// Find and remove the telegram_token field
		var newFields core.FieldsList
		for _, field := range usersCollection.Fields {
			if field.GetName() != "telegram_token" {
				newFields = append(newFields, field)
			}
		}
		
		usersCollection.Fields = newFields
		return app.Save(usersCollection)
	})
}