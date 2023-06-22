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
					"Alive",
					"Dead",
					"Pending"
				]
			}
		}`), edit_status)
		collection.Schema.AddField(edit_status)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("rpw6ibkds3yu215")
		if err != nil {
			return err
		}

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

		return dao.SaveCollection(collection)
	})
}
