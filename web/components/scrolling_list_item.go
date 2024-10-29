package wasmcomponents

import "github.com/maxence-charriere/go-app/v10/pkg/app"

// ScrollingListItem is a wasm component that renders an item in a scrolling list.
// This component uses Bootstrap.
type ScrollingListItem struct {
	app.Compo

	// Avatar is the component that renders the avatar of the item.
	// This is optional.
	Avatar app.UI

	// Title is the title of the item.
	Title string

	// Subtitle is the subtitle of the item.
	Subtitle string

	// Actions are the actions that can be performed on the item.
	// This is optional.
	Actions app.UI
}

// Render renders the item.
func (i ScrollingListItem) Render() app.UI {
	return app.Div().Class("d-flex").Body(
		app.If(i.Avatar != nil, func() app.UI {
			return app.Div().Class("flex-shrink-0").Body(i.Avatar)
		}),
		app.Div().Class("flex-grow-1 ms-3").Body(
			app.H5().Class("mb-0").Body(
				app.Text(i.Title),
			),
			app.Small().Class("text-muted").Body(
				app.Text(i.Subtitle),
			),
		),
		app.If(i.Actions != nil, func() app.UI {
			return app.Div().Class("ms-3").Body(i.Actions)
		}),
	)
}
