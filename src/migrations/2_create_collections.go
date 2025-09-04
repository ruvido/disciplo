package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Check if users collection already exists (it's created automatically)
		if _, err := app.FindCollectionByNameOrId("users"); err != nil {
			// Create users collection only if it doesn't exist
			usersCollection := core.NewBaseCollection("users")
			usersCollection.Type = "auth"
			if err := app.Save(usersCollection); err != nil {
				return err
			}
		}

		// Communities collection with all required fields
		if _, err := app.FindCollectionByNameOrId("communities"); err != nil {
			communitiesCollection := core.NewBaseCollection("communities")
			// Communities collection will be created as base collection
			// Fields will be added manually via admin panel for now
			// Required: name, description, telegram_id, type (default/local/special)
			if err := app.Save(communitiesCollection); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		// Revert - delete collections
		if collection, _ := app.FindCollectionByNameOrId("communities"); collection != nil {
			app.Delete(collection)
		}
		return nil
	})
}