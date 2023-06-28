package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

type CheckUsernameRequest struct {
	Username string `json:"username"`
}

type CheckEmailRequest struct {
	Email string `json:"email"`
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
			Path:   "/api/checkusername",
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
			Path:   "/api/checkemail",
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

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
