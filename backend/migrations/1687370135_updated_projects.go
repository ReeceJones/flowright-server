package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("rpw6ibkds3yu215")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("h8k8xdcf")

		// update
		edit_owner := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "bkx9yclv",
			"name": "owner",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"collectionId": "_pb_users_auth_",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": []
			}
		}`), edit_owner)
		collection.Schema.AddField(edit_owner)

		// update
		edit_status := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "jtrpdbva",
			"name": "status",
			"type": "select",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"alive",
					"dead",
					"pending"
				]
			}
		}`), edit_status)
		collection.Schema.AddField(edit_status)

		// update
		edit_name := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "qpl88jsn",
			"name": "name",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": 1,
				"max": 100,
				"pattern": "^[a-zA-Z]+]$"
			}
		}`), edit_name)
		collection.Schema.AddField(edit_name)

		// update
		edit_long_name := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "csnmnscs",
			"name": "long_name",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": 1,
				"max": 120,
				"pattern": "^[a-zA-Z0-9_\\-, ;?]+$"
			}
		}`), edit_long_name)
		collection.Schema.AddField(edit_long_name)

		// update
		edit_description := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "6rtjymmy",
			"name": "description",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": 1,
				"max": 200,
				"pattern": ""
			}
		}`), edit_description)
		collection.Schema.AddField(edit_description)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("rpw6ibkds3yu215")
		if err != nil {
			return err
		}

		// add
		del_slug := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "h8k8xdcf",
			"name": "slug",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), del_slug)
		collection.Schema.AddField(del_slug)

		// update
		edit_owner := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "bkx9yclv",
			"name": "owner",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"collectionId": "_pb_users_auth_",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": []
			}
		}`), edit_owner)
		collection.Schema.AddField(edit_owner)

		// update
		edit_status := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "jtrpdbva",
			"name": "status",
			"type": "select",
			"required": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"alive",
					"dead",
					"pending"
				]
			}
		}`), edit_status)
		collection.Schema.AddField(edit_status)

		// update
		edit_name := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "qpl88jsn",
			"name": "name",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_name)
		collection.Schema.AddField(edit_name)

		// update
		edit_long_name := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "csnmnscs",
			"name": "long_name",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_long_name)
		collection.Schema.AddField(edit_long_name)

		// update
		edit_description := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "6rtjymmy",
			"name": "description",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_description)
		collection.Schema.AddField(edit_description)

		return dao.SaveCollection(collection)
	})
}
