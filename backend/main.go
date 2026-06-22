package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vn4x/meswin/internal/api"
	"github.com/vn4x/meswin/internal/cv"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "github.com/vn4x/meswin/migrations"
)

func main() {
	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		TemplateLang: migratecmd.TemplateLangGo,
		Automigrate:  true,
		Dir:          "./migrations",
	})

	cvEngine := cv.NewStubEngine()
	handler := api.NewHandler(app, cvEngine)

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		handler.Register(e)

		publicDir := os.Getenv("MESWIN_PUBLIC_DIR")
		if publicDir == "" {
			publicDir = "./pb_public"
		}
		if info, err := os.Stat(publicDir); err == nil && info.IsDir() {
			e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(publicDir), true))
		}

		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/meswin/health",
			Handler: func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]string{
					"status":  "ok",
					"service": "meswin",
				})
			},
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
