package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("u0s9j6pdsqsuv4k")
		if err != nil {
			return err
		}

		collection.Name = "auth_link_requests"

		json.Unmarshal([]byte(`[
			"CREATE UNIQUE INDEX ` + "`" + `idx_SZmzzHD` + "`" + ` ON ` + "`" + `auth_link_requests` + "`" + ` (` + "`" + `challenge` + "`" + `)"
		]`), &collection.Indexes)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("u0s9j6pdsqsuv4k")
		if err != nil {
			return err
		}

		collection.Name = "auth_link_request"

		json.Unmarshal([]byte(`[
			"CREATE UNIQUE INDEX ` + "`" + `idx_SZmzzHD` + "`" + ` ON ` + "`" + `auth_link_request` + "`" + ` (` + "`" + `challenge` + "`" + `)"
		]`), &collection.Indexes)

		return dao.SaveCollection(collection)
	})
}
