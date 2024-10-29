package web

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"

	wasmcomponents "github.com/xfrr/finantrack/web/components"
	websections "github.com/xfrr/finantrack/web/sections"
)

const (
	MainPath = "/"

	AssetsPath   = "/assets"
	ExpensesPath = "/expenses"
	IncomePath   = "/income"
	ReportsPath  = "/reports"
	SettingsPath = "/settings"
)

type App struct {
	app.Compo
}

func (a *App) getSection() app.UI {
	currentPath := app.Window().URL().Path

	switch currentPath {
	case MainPath:
		return &websections.HomeSection{}
	case AssetsPath:
		return &websections.AssetsSection{
			Assets: []*websections.Asset{
				{
					Name:        "Bank",
					Description: "Bank account",
					Type:        "Bank",
					Money:       1000.0,
				},
			},
		}
	default:
		return app.Text("Not found")
	}
}

func (a *App) Render() app.UI {
	return app.Div().
		Body(
			&wasmcomponents.Navbar{
				LogoPath: "/web/static/img/logo.png",
				Title:    "Finantrack",
				Subtitle: "Your personal finance tracker",
			},
			app.Main().
				Class("container-fluid py-4 px-4 h-100 w-100").
				Body(a.getSection()),
		)
}

func NewHandler() *app.Handler {
	return &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
		Icon: app.Icon{
			Default: "/web/static/img/logo.png",
		},
		Styles: []string{
			// Bootstrap CSS
			"https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css",
			// Font Awesome CSS
			"https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.4/css/all.min.css",
		},
		Scripts: []string{
			"https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.min.js",
			"https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.4/js/all.min.js",
			"https://cdn.jsdelivr.net/npm/chart.js",
		},
	}
}
