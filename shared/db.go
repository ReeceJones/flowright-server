package shared

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/cockroachdb/pebble"
	"github.com/dchest/uniuri"
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
var pebble_db *pebble.DB = nil

func Init() (*gorm.DB, *pebble.DB) {
	db_handle, err := gorm.Open(sqlite.Open("../test.db"), &gorm.Config{})
	if err != nil {
		println(err.Error())
		panic("Failed to open database connection.")
	}

	db_handle.AutoMigrate(&BackendRef{})
	db_handle.AutoMigrate(&RoutingRule{})

	db = db_handle

	pebble_db_handle, err := pebble.Open("../test.pebble", &pebble.Options{})
	if err != nil {
		println(err.Error())
		panic("Failed to open pebble database connection.")
	}

	pebble_db = pebble_db_handle

	return db_handle, pebble_db_handle
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

func CreateEnvironment(owner string, project string, requirements string, sourceTarBall []byte, pebble_persist bool) error {
	// TODO: migrate to Podman Go Bindings
	if pebble_db == nil {
		panic("Pebble database not initialized.")
	}

	// make sure the container is present
	err := exec.Command("podman", "pull", "python:3.11").Run()
	if err != nil {
		log.Println("Error pulling python:3.11 image:", err)
		return err
	}

	containerName := uniuri.NewLenChars(10, []byte("abcdefghijklmnopqrstuvxyz123456789")) + "-flowright-" + owner + "-" + project
	buildImageName := containerName + "-build"

	defer func() {
		exec.Command("podman", "stop", containerName).Run()
		exec.Command("podman", "rm", "-f", containerName).Run()
	}()

	// create a container
	if err := exec.Command("podman", "run", "-d", "--name", containerName, "python:3.11", "sleep", "180").Run(); err != nil {
		log.Println("Error creating container:", err)
		return err
	}

	// install requirements
	tmpRequirementsFile := "/tmp/" + containerName + "_requirements.txt"
	if err := os.WriteFile(tmpRequirementsFile, []byte(requirements), 0644); err != nil {
		log.Println("Error writing requirements file:", err)
		return err
	}

	if err := exec.Command("podman", "cp", tmpRequirementsFile, containerName+":/requirements.txt").Run(); err != nil {
		log.Println("Error copying requirements file:", err)
		return err
	}

	if err := exec.Command("podman", "exec", containerName, "pip", "install", "-r", "/requirements.txt").Run(); err != nil {
		log.Println("Error installing requirements:", err)
		return err
	}

	// ensure flowright is installed (this issue mostly occurs with flowright development)
	if err := exec.Command("podman", "exec", containerName, "pip", "install", "flowright").Run(); err != nil {
		log.Println("Error installing flowright:", err)
		return err
	}

	// copy source
	tmpSourceFile := "/tmp/" + containerName + "_source.tar.gz"
	if err := os.WriteFile(tmpSourceFile, sourceTarBall, 0644); err != nil {
		log.Println("Error writing source file:", err)
		return err
	}

	if err := exec.Command("podman", "cp", tmpSourceFile, containerName+":/source.tar.gz").Run(); err != nil {
		log.Println("Error copying source:", err)
		return err
	}

	if err := exec.Command("podman", "exec", containerName, "mkdir", "-p", "/flowright_app").Run(); err != nil {
		log.Println("Error creating app directory:", err)
		return err
	}

	if err := exec.Command("podman", "exec", containerName, "tar", "-xzf", "/source.tar.gz", "-C", "/flowright_app").Run(); err != nil {
		log.Println("Failed to extract source:", err)
		return err
	}

	// commit changes to image
	commitCommand := exec.Command("podman", "commit", containerName, buildImageName)
	// commitCommand.Stdout = os.Stdout // TODO: capture this output
	// commitCommand.Stderr = os.Stderr

	if err := commitCommand.Run(); err != nil {
		log.Println("Error committing changes:", err)
		return err
	}

	if pebble_persist {
		// save to pebble db
		imageWriter := new(bytes.Buffer)
		saveCommand := exec.Command("podman", "save", "--compress", buildImageName)
		saveCommand.Stdout = imageWriter

		if err := saveCommand.Run(); err != nil {
			log.Println("Error saving image:", err)
			return err
		}

		if err := pebble_db.Set([]byte(buildImageName), imageWriter.Bytes(), pebble.Sync); err != nil {
			log.Println("Error saving image to pebble db:", err)
			return err
		}
	}

	// stop base container
	if err := exec.Command("podman", "stop", containerName).Run(); err != nil {
		log.Println("Error stopping container:", err)
		return err
	}

	// now run app
	// TODO: expose volume with UDS
	containerRunName := containerName + "-run"
	runCommand := exec.Command("podman", "run", "-d", "-p", "8000:8000", "--name", containerRunName, buildImageName, "flowright", "run", "/flowright_app", "--host=0.0.0.0")
	runCommand.Stderr = os.Stderr
	runCommand.Stdout = os.Stdout
	if err := runCommand.Run(); err != nil {
		log.Println("Error running app:", err)
		return err
	}

	return nil
}
