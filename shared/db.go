package shared

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BackendRef struct {
	gorm.Model
	Owner         string
	Project       string
	RoutingRuleID uint
	Route         RoutingRule
}

type RoutingRule struct {
	gorm.Model
	BackendRefID uint
	Endpoint     string
}

var db *gorm.DB = nil

func Init() *gorm.DB {
	db_handle, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to open database connection.")
	}

	db_handle.AutoMigrate(&BackendRef{})
	db_handle.AutoMigrate(&RoutingRule{})

	db = db_handle

	return db_handle
}

func AddRoutingRule(owner string, project string, endpoint string) {
	existing_route, err := GetRoute(owner, project)
	if err != nil || existing_route == nil {
		ref := BackendRef{Owner: owner, Project: project, Route: RoutingRule{Endpoint: endpoint}}
		db.Save(&ref)
	} else {
		existing_route.Endpoint = endpoint
		db.Save(&existing_route.Endpoint)
	}
}

func GetRoute(owner string, project string) (route *RoutingRule, err error) {
	ref := &BackendRef{}
	// result := db.Model(&BackendRef{}).Where(&BackendRef{Owner: owner, Project: project}).Take(&ref)
	result := db.Preload("Route").Where(&BackendRef{Owner: owner, Project: project}).Take(&ref)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ref.Route, nil
}
