package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Add telegram_token_created field to users collection for token expiration
		usersCollection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		// Add telegram_token_created field (datetime)
		telegramTokenCreatedField := &core.DateField{
			Id:   "telegram_token_created",
			Name: "telegram_token_created",
		}
		
		usersCollection.Fields = append(usersCollection.Fields, telegramTokenCreatedField)

		return app.Save(usersCollection)
	}, func(app core.App) error {
		// Revert - remove the field
		usersCollection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		// Remove telegram_token_created field
		for i, field := range usersCollection.Fields {
			if field.GetId() == "telegram_token_created" {
				usersCollection.Fields = append(usersCollection.Fields[:i], usersCollection.Fields[i+1:]...)
				break
			}
		}

		return app.Save(usersCollection)
	})
}