package wasmcomponents

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Navbar struct {
	app.Compo

	LogoPath string
	Title    string
	Subtitle string

	currentPath string
}

func (n *Navbar) OnNav(ctx app.Context) {
	path := ctx.Page().URL().Path
	if path != "" && path != n.currentPath {
		n.currentPath = path
	}
}

func newNavItemLink(text, href string, active bool) app.UI {
	return &ItemLinkComponent{
		Text:   text,
		Href:   href,
		Active: active,
	}
}

func (n *Navbar) Render() app.UI {
	isActive := func(path string) bool {
		return n.currentPath == path
	}

	return app.Nav().
		Class("navbar navbar-expand-lg navbar-dark").
		Style("background-color", "#343a40").
		Body(
			app.A().Href("/").Class("navbar-brand").Body(
				app.Img().
					Width(60).Height(60).
					Class("navbar-brand-img rounded-circle").
					Src(n.LogoPath).Alt(n.Title),
				app.Text(n.Title),
			),
			app.
				Button().
				Class("navbar-toggler").
				DataSet("bs-toggle", "collapse").
				DataSet("bs-target", "#navbarNav").
				Body(
					app.Span().Class("navbar-toggler-icon"),
				),
			app.Div().Class("collapse navbar-collapse").ID("navbarNav").Body(
				app.Ul().Class("navbar-nav ms-auto").Body(
					newNavItemLink("Home", "/", isActive("/")),
					newNavItemLink("Assets", "/assets", isActive("/assets")),
					newNavItemLink("Expenses", "/expenses", isActive("/expenses")),
					newNavItemLink("Reports", "/reports", isActive("/reports")),
					newNavItemLink("Settings", "/settings", isActive("/settings")),
				),
			),
		)
}
