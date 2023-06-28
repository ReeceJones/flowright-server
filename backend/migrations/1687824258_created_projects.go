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
			"id": "u6my0iagdcpi3id",
			"created": "2023-06-27 00:04:18.055Z",
			"updated": "2023-06-27 00:04:18.055Z",
			"name": "projects",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "msgk4gk9",
					"name": "name",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": 1,
						"max": 100,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "4hfv03gp",
					"name": "description",
					"type": "text",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "oh0jloop",
					"name": "visibility",
					"type": "select",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"values": [
							"private",
							"public"
						]
					}
				},
				{
					"system": false,
					"id": "fxlgdl5c",
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
				},
				{
					"system": false,
					"id": "vhqs0h1b",
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
				}
			],
			"indexes": [],
			"listRule": "@request.auth.id = owner.id",
			"viewRule": "@request.auth.id = owner.id",
			"createRule": "@request.auth.id = owner.id",
			"updateRule": "@request.auth.id = owner.id",
			"deleteRule": "@request.auth.id = owner.id",
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("u6my0iagdcpi3id")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
