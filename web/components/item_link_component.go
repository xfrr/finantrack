package wasmcomponents

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ItemLinkComponent is a wasm component that renders a link to an item.
type ItemLinkComponent struct {
	app.Compo

	Text   string
	Href   string
	Active bool
}

// Render renders the item link component.
func (i *ItemLinkComponent) Render() app.UI {
	activeClass := ""
	if i.Active {
		activeClass = "active"
	}

	return app.Li().
		Class("nav-item").
		Body(
			app.A().Href(i.Href).Class("nav-link " + activeClass).Body(
				app.Text(i.Text),
			),
		)
}
