package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "u0s9j6pdsqsuv4k",
			"created": "2023-06-29 02:14:12.492Z",
			"updated": "2023-06-29 02:14:12.492Z",
			"name": "auth_link_request",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "lk8jmdoo",
					"name": "challenge",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": 64,
						"max": 64,
						"pattern": "^[a-zA-Z0-9]{64}$"
					}
				},
				{
					"system": false,
					"id": "kselsdbq",
					"name": "success",
					"type": "bool",
					"required": false,
					"unique": false,
					"options": {}
				},
				{
					"system": false,
					"id": "yyakyvdt",
					"name": "current_client_jwt",
					"type": "text",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				}
			],
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_SZmzzHD` + "`" + ` ON ` + "`" + `auth_link_request` + "`" + ` (` + "`" + `challenge` + "`" + `)"
			],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("u0s9j6pdsqsuv4k")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
