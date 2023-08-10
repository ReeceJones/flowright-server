package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "flowright/backend/reece.ooo/flowright"
)

var (
	controlPlaneHost = os.Getenv("CONTROL_PLANE_HOST")
	controlPlanePort = os.Getenv("CONTROL_PLANE_PORT")
)

type CheckUsernameRequest struct {
	Username string `json:"username"`
}

type CheckEmailRequest struct {
	Email string `json:"email"`
}

type RequirementEntry struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type UploadProjectRequest struct {
	ProjectId    string             `json:"project_id"`
	Data         string             `json:"data"`
	Requirements []RequirementEntry `json:"requirements"`
}

func main() {
	app := pocketbase.New()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		// enable auto creation of migration files when making collection changes
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/flowright/checkusername",
			Handler: func(c echo.Context) error {
				userRequest := new(CheckUsernameRequest)
				if err := c.Bind(userRequest); err != nil {
					return err
				}

				matchingUsers := 0
				err := app.DB().
					Select("count(*)").
					From("users").
					OrWhere(dbx.HashExp{"username": userRequest.Username}).
					Row(&matchingUsers)

				if err != nil {
					return err
				}

				if matchingUsers > 0 {
					return c.String(http.StatusConflict, "Username exists already.")
				}

				return c.String(http.StatusAccepted, "ok")
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app), apis.RequireGuestOnly(),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/flowright/checkemail",
			Handler: func(c echo.Context) error {
				userRequest := new(CheckEmailRequest)
				if err := c.Bind(userRequest); err != nil {
					return err
				}

				matchingUsers := 0
				err := app.DB().
					Select("count(*)").
					From("users").
					OrWhere(dbx.HashExp{"email": userRequest.Email}).
					Row(&matchingUsers)

				if err != nil {
					return err
				}

				if matchingUsers > 0 {
					return c.String(http.StatusConflict, "Email is already in use.")
				}

				return c.String(http.StatusAccepted, "ok")
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app), apis.RequireGuestOnly(),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/flowright/auth_link",
			Handler: func(c echo.Context) error {
				collection, err := app.Dao().FindCollectionByNameOrId("auth_link_requests")
				if err != nil {
					return err
				}

				record := models.NewRecord(collection)

				if err := app.Dao().SaveRecord(record); err != nil {
					return err
				}

				data := map[string]string{
					"id": record.Id,
				}

				return c.JSON(http.StatusCreated, data)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app), apis.RequireGuestOnly(),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/flowright/auth_link/:id",
			Handler: func(c echo.Context) error {
				record, err := app.Dao().FindRecordById("auth_link_requests", c.PathParam("id"))
				if err != nil {
					return err
				}

				if !record.GetBool("success") {
					return c.String(http.StatusTeapot, "Not linked yet.")
				}

				app.Dao().DeleteRecord(record)

				return c.JSON(http.StatusAccepted, record)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app), apis.RequireGuestOnly(),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/flowright/upload",
			Handler: func(c echo.Context) error {
				uploadRequest := new(UploadProjectRequest)
				if err := c.Bind(uploadRequest); err != nil {
					return err
				}

				project, err := app.Dao().FindRecordById("projects", uploadRequest.ProjectId)
				if err != nil {
					return err
				}

				owner, err := app.Dao().FindRecordById("users", project.GetString("owner"))
				if err != nil {
					return err
				}

				collection, err := app.Dao().FindCollectionByNameOrId("project_uploads")
				if err != nil {
					return err
				}

				tarData, err := base64.StdEncoding.DecodeString(uploadRequest.Data)
				if err != nil {
					return err
				}

				record := models.NewRecord(collection)
				record.Set("project", uploadRequest.ProjectId)
				record.Set("data", tarData)
				record.Set("requirements", uploadRequest.Requirements)

				if err := app.Dao().SaveRecord(record); err != nil {
					return err
				}

				requirements := ""
				for _, requirement := range uploadRequest.Requirements {
					requirements += requirement.Name + "==" + requirement.Version + "\n"
				}

				// create base container state
				// TODO: should have configuratble control plane addr
				conn, err := grpc.Dial(fmt.Sprintf("%s:%s", controlPlaneHost, controlPlanePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Print(err)
					log.Print("Failed to connect to proxy")
					return err
				}
				defer conn.Close()

				client := pb.NewRoutingControllerClient(conn)
				cctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
				defer cancel()

				environmentResponse, err := client.CreateEnvironment(cctx, &pb.EnvironmentCreateRequest{
					Owner:          owner.Username(),
					Project:        project.GetString("name"),
					Requirements:   requirements,
					ProjectTarball: tarData,
				})

				// err = CreateEnvironment(owner.Username(), project.GetString("name"), requirements, tarData)
				if err != nil {
					log.Printf("Failed to create environment (%s, %s)\n", owner.Username(), project.GetString("name"))
					return c.String(http.StatusInternalServerError, err.Error())
				}

				// create routing rule
				_, err = client.CreateOrUpdateRoute(cctx, &pb.RoutingMap{
					Owner:     owner.Username(),
					Project:   project.GetString("name"),
					ProxyName: environmentResponse.ProxxyName, // FIXME: oops, typo
				})

				if err != nil {
					log.Printf("Failed to create route for environment (%s, %s) on %s\n", owner.Username(), project.GetString("name"), environmentResponse.ProxxyName)
					return err
				}

				data := map[string]string{
					"id": record.Id,
				}

				return c.JSON(http.StatusCreated, data)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app), apis.RequireAdminOrRecordAuth("users"),
			},
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
