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
			"id": "84nzxzjd2env3yz",
			"created": "2023-07-04 00:21:05.526Z",
			"updated": "2023-07-04 00:21:05.526Z",
			"name": "project_uploads",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "mxytumjq",
					"name": "project",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"collectionId": "u6my0iagdcpi3id",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "jdlwqbox",
					"name": "requirements",
					"type": "json",
					"required": false,
					"unique": false,
					"options": {}
				},
				{
					"system": false,
					"id": "2rwgnorp",
					"name": "data",
					"type": "file",
					"required": false,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"maxSize": 5242880,
						"mimeTypes": [
							"application/x-tar"
						],
						"thumbs": [],
						"protected": false
					}
				}
			],
			"indexes": [],
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

		collection, err := dao.FindCollectionByNameOrId("84nzxzjd2env3yz")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
