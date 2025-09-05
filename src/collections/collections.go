package collections

import (
	"fmt"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Initialize(app *pocketbase.PocketBase) error {
	if err := createUsersCollection(app); err != nil {
		return err
	}
	
	if err := createCommunitiesCollection(app); err != nil {
		return err
	}
	
	if err := createTokensCollection(app); err != nil {
		return err
	}
	
	if err := createRequestsCollection(app); err != nil {
		return err
	}
	
	return nil
}

func createUsersCollection(app *pocketbase.PocketBase) error {
	return app.Dao().RunInTransaction(func(txDao *core.DAO) error {
		collection, _ := txDao.FindCollectionByNameOrId("users")
		if collection != nil {
			return nil
		}

		collection = &core.Collection{
			Name:       "users",
			Type:       core.CollectionTypeBase,
			ListRule:   core.NewString("@request.auth.admin = true"),
			ViewRule:   core.NewString("@request.auth.id = id || @request.auth.admin = true"),
			CreateRule: core.NewString("@request.auth.admin = true"),
			UpdateRule: core.NewString("@request.auth.id = id || @request.auth.admin = true"),
			DeleteRule: core.NewString("@request.auth.admin = true"),
		}

		collection.Schema.AddField(&core.SchemaField{
			Name:     "name",
			Type:     core.FieldTypeText,
			Required: true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name:     "email",
			Type:     core.FieldTypeEmail,
			Required: true,
			Unique:   true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name:     "password",
			Type:     core.FieldTypeText,
			Required: false,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "metadata",
			Type: core.FieldTypeJson,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "admin",
			Type: core.FieldTypeBool,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "group_admin",
			Type: core.FieldTypeRelation,
			Options: &core.RelationOptions{
				CollectionId:  "communities",
				CascadeDelete: false,
				MinSelect:     core.NewInt(0),
				MaxSelect:     core.NewInt(999),
			},
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "group_admin_since",
			Type: core.FieldTypeDate,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "groups",
			Type: core.FieldTypeRelation,
			Options: &core.RelationOptions{
				CollectionId:  "communities",
				CascadeDelete: false,
				MinSelect:     core.NewInt(0),
				MaxSelect:     core.NewInt(999),
			},
		})

		collection.Schema.AddField(&core.SchemaField{
			Name:    "status",
			Type:    core.FieldTypeSelect,
			Options: &core.SelectOptions{
				Values: []string{"pending", "accepted"},
			},
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "verified",
			Type: core.FieldTypeBool,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "telegram_id",
			Type: core.FieldTypeNumber,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "telegram_name",
			Type: core.FieldTypeText,
		})

		return txDao.SaveCollection(collection)
	})
}

func createCommunitiesCollection(app *pocketbase.PocketBase) error {
	return app.Dao().RunInTransaction(func(txDao *core.DAO) error {
		collection, _ := txDao.FindCollectionByNameOrId("communities")
		if collection != nil {
			return nil
		}

		collection = &core.Collection{
			Name:       "communities",
			Type:       core.CollectionTypeBase,
			ListRule:   core.NewString("@request.auth.admin = true"),
			ViewRule:   core.NewString("@request.auth.admin = true"),
			CreateRule: core.NewString("@request.auth.admin = true"),
			UpdateRule: core.NewString("@request.auth.admin = true"),
			DeleteRule: core.NewString("@request.auth.admin = true"),
		}

		collection.Schema.AddField(&core.SchemaField{
			Name:     "name",
			Type:     core.FieldTypeText,
			Required: true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "description",
			Type: core.FieldTypeText,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "telegram_chat_id",
			Type: core.FieldTypeNumber,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "telegram_username",
			Type: core.FieldTypeText,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name:    "type",
			Type:    core.FieldTypeSelect,
			Options: &core.SelectOptions{
				Values: []string{"default", "local", "special"},
			},
		})

		return txDao.SaveCollection(collection)
	})
}

func createTokensCollection(app *pocketbase.PocketBase) error {
	return app.Dao().RunInTransaction(func(txDao *core.DAO) error {
		collection, _ := txDao.FindCollectionByNameOrId("tokens")
		if collection != nil {
			return nil
		}

		collection = &core.Collection{
			Name:       "tokens",
			Type:       core.CollectionTypeBase,
			ListRule:   core.NewString("@request.auth.admin = true"),
			ViewRule:   core.NewString("@request.auth.admin = true"),
			CreateRule: core.NewString("@request.auth.admin = true"),
			UpdateRule: core.NewString("@request.auth.admin = true"),
			DeleteRule: core.NewString("@request.auth.admin = true"),
		}

		collection.Schema.AddField(&core.SchemaField{
			Name:     "email",
			Type:     core.FieldTypeEmail,
			Required: true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name:     "token",
			Type:     core.FieldTypeText,
			Required: true,
			Unique:   true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name:    "type",
			Type:    core.FieldTypeSelect,
			Options: &core.SelectOptions{
				Values: []string{"admin_link", "user_invite"},
			},
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "used",
			Type: core.FieldTypeBool,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "expires_at",
			Type: core.FieldTypeDate,
		})

		return txDao.SaveCollection(collection)
	})
}

func createRequestsCollection(app *pocketbase.PocketBase) error {
	return app.Dao().RunInTransaction(func(txDao *core.DAO) error {
		collection, _ := txDao.FindCollectionByNameOrId("requests")
		if collection != nil {
			return nil
		}

		collection = &core.Collection{
			Name:       "requests",
			Type:       core.CollectionTypeBase,
			ListRule:   core.NewString("@request.auth.admin = true"),
			ViewRule:   core.NewString("@request.auth.admin = true"),
			CreateRule: core.NewString(""),  // Allow public registration
			UpdateRule: core.NewString("@request.auth.admin = true"),
			DeleteRule: core.NewString("@request.auth.admin = true"),
		}

		// Basic registration data
		collection.Schema.AddField(&core.SchemaField{
			Name:     "name",
			Type:     core.FieldTypeText,
			Required: true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name:     "email",
			Type:     core.FieldTypeEmail,
			Required: true,
			Unique:   true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name:     "password",
			Type:     core.FieldTypeText,
			Required: true,
		})

		// Personal details
		collection.Schema.AddField(&core.SchemaField{
			Name:     "date_of_birth",
			Type:     core.FieldTypeDate,
			Required: true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name:     "city",
			Type:     core.FieldTypeText,
			Required: true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name:     "location",
			Type:     core.FieldTypeText,
			Required: true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name:     "job_field",
			Type:     core.FieldTypeText,
			Required: true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "interests",
			Type: core.FieldTypeJson,
			Required: true,
		})

		// Application details
		collection.Schema.AddField(&core.SchemaField{
			Name:     "why_join",
			Type:     core.FieldTypeText,
			Required: true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "profile_picture",
			Type: core.FieldTypeFile,
			Options: &core.FileOptions{
				MaxSelect: 1,
				MaxSize:   5242880, // 5MB
				MimeTypes: []string{"image/jpeg", "image/jpg", "image/png"},
			},
			Required: true,
		})

		// Status and approval tracking
		collection.Schema.AddField(&core.SchemaField{
			Name:    "status",
			Type:    core.FieldTypeSelect,
			Options: &core.SelectOptions{
				Values: []string{"pending", "approved", "rejected"},
			},
			Required: true,
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "approved_by",
			Type: core.FieldTypeRelation,
			Options: &core.RelationOptions{
				CollectionId:  "users",
				CascadeDelete: false,
			},
		})

		collection.Schema.AddField(&core.SchemaField{
			Name: "approved_at",
			Type: core.FieldTypeDate,
		})

		// User creation reference (when approved)
		collection.Schema.AddField(&core.SchemaField{
			Name: "created_user_id",
			Type: core.FieldTypeRelation,
			Options: &core.RelationOptions{
				CollectionId:  "users",
				CascadeDelete: false,
			},
		})

		return txDao.SaveCollection(collection)
	})
}