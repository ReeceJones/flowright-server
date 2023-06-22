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
				"pattern": "^[a-zA-Z]+$"
			}
		}`), edit_name)
		collection.Schema.AddField(edit_name)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("rpw6ibkds3yu215")
		if err != nil {
			return err
		}

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

		return dao.SaveCollection(collection)
	})
}
