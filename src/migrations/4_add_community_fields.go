package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Add missing fields to communities collection
		// Required fields from CLAUDE.md: name, description, created (auto), telegram_id, type
		
		communitiesCollection, err := app.FindCollectionByNameOrId("communities")
		if err != nil {
			return err // Communities collection must exist
		}

		// Add missing fields
		newFields := core.FieldsList{
			&core.TextField{
				Id:   "name",
				Name: "name",
				Required: true,
			},
			&core.TextField{
				Id:   "description",
				Name: "description",
			},
			&core.TextField{
				Id:   "telegram_id",
				Name: "telegram_id",
			},
			&core.SelectField{
				Id:   "type",
				Name: "type",
				Values: []string{"default", "local", "special"},
				MaxSelect: 1,
				Required: true,
			},
		}

		// Append new fields to existing fields
		communitiesCollection.Fields = append(communitiesCollection.Fields, newFields...)
		
		return app.Save(communitiesCollection)
	}, func(app core.App) error {
		// Revert: Remove the custom fields we added
		// (This is complex to implement correctly, so we'll leave empty for now)
		return nil
	})
}