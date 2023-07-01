package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("u0s9j6pdsqsuv4k")
		if err != nil {
			return err
		}

		collection.ListRule = nil

		collection.ViewRule = types.Pointer("")

		json.Unmarshal([]byte(`[]`), &collection.Indexes)

		// remove
		collection.Schema.RemoveField("lk8jmdoo")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("u0s9j6pdsqsuv4k")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("")

		collection.ViewRule = nil

		json.Unmarshal([]byte(`[
			"CREATE UNIQUE INDEX ` + "`" + `idx_SZmzzHD` + "`" + ` ON ` + "`" + `auth_link_requests` + "`" + ` (` + "`" + `challenge` + "`" + `)"
		]`), &collection.Indexes)

		// add
		del_challenge := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), del_challenge)
		collection.Schema.AddField(del_challenge)

		return dao.SaveCollection(collection)
	})
}
