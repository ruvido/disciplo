package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Add remaining missing fields to users collection
		// Still need: group_admin (relation), group_admin_since (date), 
		// groups (relation), status (select)
		
		usersCollection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err // Users collection must exist
		}

		// Get communities collection ID
		communitiesCollection, err := app.FindCollectionByNameOrId("communities")
		if err != nil {
			return err // Communities collection must exist for relations
		}

		// Add remaining missing fields
		newFields := core.FieldsList{
			&core.RelationField{
				Id:   "group_admin",
				Name: "group_admin",
				CollectionId: communitiesCollection.Id, // Use actual collection ID
			},
			&core.DateField{
				Id:   "group_admin_since",
				Name: "group_admin_since",
			},
			&core.RelationField{
				Id:   "groups",
				Name: "groups", 
				CollectionId: communitiesCollection.Id, // Use actual collection ID
			},
			&core.SelectField{
				Id:   "status",
				Name: "status",
				Values: []string{"pending", "accepted"},
				MaxSelect: 1,
				Required: true,
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